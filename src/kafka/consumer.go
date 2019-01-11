package kafka

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"siren/configs"
	"siren/models"
	"siren/pkg/database"
	"siren/pkg/utils"
	"strings"
	"time"

	"github.com/tidwall/gjson"

	"siren/src/workers"

	"github.com/bsm/sarama-cluster"
)

type CountFrequentConsumerParamsType struct {
	brokers []string
	groupID string
	topics  []string
	handler func([]byte)
}

var MallCountFrequentConsumerParams CountFrequentConsumerParamsType
var StoreCountFrequentConsumerParams CountFrequentConsumerParamsType

var titanParams struct {
	identificationURL string
	groupCreateURL    string
}

var ruleNumber struct {
	high int
	low  int
}

func consumerInit() {
	// todo: fixed use config.fetchValue
	host := configs.FetchFieldValue("KAFKAHOST")
	port := configs.FetchFieldValue("KAFKAPORT")
	groupName := configs.FetchFieldValue("KAFKAGROUP")
	topic := configs.FetchFieldValue("KAFKATOPIC")
	log.Println(fmt.Sprintf("env: %s,host:%s, port: %s, groupID: %s, topic: %s", configs.ENV, host, port, groupName, topic))
	MallCountFrequentConsumerParams.brokers = []string{fmt.Sprintf("%s:%s", host, port)}
	MallCountFrequentConsumerParams.groupID = groupName
	MallCountFrequentConsumerParams.topics = []string{topic}
	MallCountFrequentConsumerParams.handler = mallInfoHandler

	StoreCountFrequentConsumerParams.brokers = []string{fmt.Sprintf("%s:%s", host, port)}
	StoreCountFrequentConsumerParams.groupID = groupName
	StoreCountFrequentConsumerParams.topics = []string{"store_frequent_customer_" + configs.ENV}
	StoreCountFrequentConsumerParams.handler = storeInfoHandler

	// todo: titan faces/identification
	titanParams.identificationURL = fmt.Sprintf(configs.FetchFieldValue("TitanHOST") + "/faces/identification")
	titanParams.groupCreateURL = fmt.Sprintf(configs.FetchFieldValue("TitanHOST") + "/groups/create")
	// todo: rule number
	ruleNumber.high = 3
	ruleNumber.low = 2
}

func CountFrequentConsumer() {
	consumerInit()
	go MallCountFrequentConsumerParams.StartConsumer()
	//StoreCountFrequentConsumerParams.StartConsumer()
}

func (params *CountFrequentConsumerParamsType) StartConsumer() {
	config := cluster.NewConfig()
	config.Consumer.Return.Errors = true
	config.Group.Return.Notifications = true
	consumer, err := cluster.NewConsumer(params.brokers, params.groupID, params.topics, config)
	if err != nil {
		panic(err)
	}
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	// consume errors
	go func() {
		for err := range consumer.Errors() {
			log.Printf("Error: %s\n", err.Error())
		}
	}()

	// consume notifications
	go func() {
		for ntf := range consumer.Notifications() {
			log.Printf("Rebalanced: %+v\n", ntf)
		}
	}()

	// consume messages, watch signals
	for {
		select {
		case msg, ok := <-consumer.Messages():
			if ok {
				fmt.Fprintf(os.Stdout, "%s/%d/%d\t%s\t%s\n", msg.Topic, msg.Partition, msg.Offset, msg.Key, msg.Value)
				params.handler(msg.Value)
				consumer.MarkOffset(msg, "") // mark message as processed
			}
		case <-signals:
			return
		}
	}
}

type InfoForKafkaProducer struct {
	CompanyID  uint   `json:"company_id"`
	FaceID     string `json:"face_id"`
	PersonID   string `json:"person_id"`
	CapturedAt int64  `json:"captured_at"`
	EventID    uint   `json:"event_id"`
}

func mallInfoHandler(values []byte) {
	var info InfoForKafkaProducer
	if err := json.Unmarshal(values, &info); err != nil {
		log.Println(err)
		return
	}

	var group *models.FrequentCustomerGroup
	var ok bool
	if ok, group = saveGroupInfo(info.CompanyID); !ok {
		return
	}
	fetchDataByTitan(group, info)

}

func saveGroupInfo(companyID uint) (bool, *models.FrequentCustomerGroup) {
	var oneGroup models.FrequentCustomerGroup
	var (
		ok      bool
		groupID string
	)
	if dbError := database.POSTGRES.Where("company_id = ?", companyID).First(&oneGroup).Error; dbError != nil {
		oneGroup = models.FrequentCustomerGroup{
			CompanyID: companyID,
		}

		if ok, groupID = titanAddGroup(oneGroup.GroupUUID, fmt.Sprintf("%d_回头客", oneGroup.CompanyID)); !ok {
			return false, nil
		}
		oneGroup.GroupUUID = groupID

		if dbError := database.POSTGRES.Save(&oneGroup).Error; dbError != nil {
			return false, nil
		}

	}
	return true, &oneGroup
}

func titanAddGroup(groupUUID string, name string) (bool, string) {
	apiID := configs.FetchFieldValue("TitanAPIID")
	apiSecret := configs.FetchFieldValue("TitanAPISecret")
	response, err := http.PostForm(titanParams.groupCreateURL, url.Values{
		"api_id":     {apiID},
		"api_secret": {apiSecret},
		"group_id":   {groupUUID},
		"name":       {name},
	})
	if err != nil {
		log.Println("create group fail")
		return false, "-1"
	}
	defer response.Body.Close()
	log.Println("titan add group", response.StatusCode)
	content, _ := ioutil.ReadAll(response.Body)
	values := gjson.ParseBytes(content)
	if values.Get("status").String() != "ok" {
		log.Println("status is not ok, create group fail")
		return false, "-1"
	}
	groupID := values.Get("group_id").String()
	return true, groupID
}

func fetchDataByTitan(group *models.FrequentCustomerGroup, info InfoForKafkaProducer) bool {
	log.Println("URL", titanParams.identificationURL)
	response, err := http.PostForm(titanParams.identificationURL, url.Values{
		"api_id":     {configs.FetchFieldValue("TitanAPIID")},
		"api_secret": {configs.FetchFieldValue("TitanAPISecret")},
		"face_id":    {info.FaceID},
		"group_id":   {group.GroupUUID},
		"top":        {"20"},
	})
	log.Println("response", response.StatusCode)
	if err != nil {
		return false
	}
	defer response.Body.Close()

	responseByte, _ := ioutil.ReadAll(response.Body)

	log.Println("titan values", responseByte)
	// todo: fix it if status is not ok
	if info.CompanyID != 0 {
		if ok := personIDHandler(info.EventID, group.ID, info.PersonID, responseByte, info.CapturedAt); !ok {
			return false
		}
	}

	return true

}

type result struct {
	PersonID string    `json:"person_id"`
	Day      time.Time `json:"day"`
}

type results []result

func personIDHandler(eventID uint, groupID uint, personUUID string, values []byte, capturedAt int64) bool {
	valuesJson := gjson.ParseBytes(values)
	exists := valuesJson.Get("candidates").Exists()
	if !exists || (exists && len(valuesJson.Get("candidates").Array()) == 0) {
		var onePerson models.FrequentCustomerPeople
		onePerson.PersonID = personUUID
		onePerson.Date = utils.CurrentDate(time.Unix(capturedAt, 0))
		hour, _ := time.Parse("2006-01-02 15:00:00", time.Unix(capturedAt, 0).Format("2006-01-02 15:00:00"))
		onePerson.Hour = hour
		onePerson.Frequency = 0
		onePerson.Interval = 0
		onePerson.FrequentCustomerGroupID = groupID
		onePerson.IsFrequentCustomer = false
		onePerson.EventID = eventID
		database.POSTGRES.Save(&onePerson)
		workers.MallCountFrequentCustomerHandler(onePerson, groupID, capturedAt)
		return true
	} else {
		var personIDs []string
		for _, i := range valuesJson.Get("candidates").Array() {
			personIDs = append(personIDs, fmt.Sprintf("`%s`", i.Get("person_id").String()))
		}
		personIDString := strings.Join(personIDs, ",")
		now := time.Now()
		right := now.Format("2006-01-02 15:04:05")
		left := now.AddDate(0, -1, 0).Format("2006-01-02 15:04:05")
		sql := fmt.Sprintf(`SELECT person_id, date_trunc('day',max(capture_at)) as day FROM events WHERE person_id in (%s) AND capture_at BETWEEN '%s' AND '%s' group by person_id ORDER BY day desc`,
			personIDString, left, right)

		var resultsValues results

		database.POSTGRES.Raw(sql).Scan(&resultsValues)

		//personID
		var onePerson models.FrequentCustomerPeople
		hour, _ := time.Parse("2006-01-02 15:00:00", time.Unix(capturedAt, 0).Format("2006-01-02 15:00:00"))
		if dbError := database.POSTGRES.Where("person_id = ? AND hour = ?", personUUID, hour).First(&onePerson).Error; dbError != nil {
			onePerson = models.FrequentCustomerPeople{
				PersonID:                personUUID,
				FrequentCustomerGroupID: groupID,
				Date:                    utils.CurrentDate(time.Unix(capturedAt, 0)),
				Hour:                    hour,
				Frequency:               uint(len(resultsValues)),
				EventID:                 eventID,
			}
			if len(resultsValues) == 0 {
				onePerson.Interval = 0 // 新客，间隔为 0
				onePerson.IsFrequentCustomer = false
			} else {
				onePerson.Interval = uint(float64(time.Now().Sub(resultsValues[0].Day).Hours()/24) + 1)
				onePerson.IsFrequentCustomer = true
			}
			if dbError := database.POSTGRES.Save(&onePerson).Error; dbError != nil {
				return false
			}
		}
		workers.MallCountFrequentCustomerHandler(onePerson, groupID, capturedAt)
	}
	return true
}

type StoreInfo struct {
	CompanyID uint   `json:"company_id"`
	ShopID    uint   `json:"shop_id"`
	PersonID  string `json:"person_id"`
	CaptureAt int64  `json:"capture_at"`
}

func storeInfoHandler(values []byte) {

	var info StoreInfo
	if err := json.Unmarshal(values, &info); err != nil {
		log.Println(err)
		return
	}

	workers.StoreFrequentCustomerHandler(info.CompanyID, info.ShopID, info.PersonID, info.CaptureAt)

}
