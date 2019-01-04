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
	"siren/pkg/titan"
	"siren/utils"
	"strings"
	"time"

	"github.com/bsm/sarama-cluster"
)

var CountFrequentConsumerParams struct {
	brokers []string
	groupID string
	topics  []string
}

var titanParams struct {
	identificationURL string
}

var ruleNumber struct {
	high int
	low  int
}

func consumerInit() {
	// todo: fix use config.fetchValue
	host := configs.FetchFieldValue("KAFKAHOST")
	port := configs.FetchFieldValue("KAFKAPORT")
	groupName := configs.FetchFieldValue("KAFKAGROUP")
	topic := configs.FetchFieldValue("KAFKATOPIC")
	log.Println(fmt.Sprintf("env: %s,host:%s, port: %s, groupID: %s, topic: %s", configs.ENV, host, port, groupName, topic))
	CountFrequentConsumerParams.brokers = []string{fmt.Sprintf("%s:%s", host, port)}
	CountFrequentConsumerParams.groupID = groupName
	CountFrequentConsumerParams.topics = []string{topic}

	// todo: titan faces/identification
	titanParams.identificationURL = fmt.Sprintf("http://" + configs.FetchFieldValue("TitanHOST") + "/faces/identification")

	// todo: rule number
	ruleNumber.high = 3
	ruleNumber.low = 2
}

func CountFrequentConsumer() {
	consumerInit()
	StartForConsumer(
		CountFrequentConsumerParams.brokers,
		CountFrequentConsumerParams.groupID,
		CountFrequentConsumerParams.topics,
	)
}

func StartForConsumer(brokers []string, groupID string, topics []string) {
	config := cluster.NewConfig()
	config.Consumer.Return.Errors = true
	config.Group.Return.Notifications = true
	consumer, err := cluster.NewConsumer(brokers, groupID, topics, config)
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
				infoHandler(msg.Value)
				consumer.MarkOffset(msg, "") // mark message as processed
			}
		case <-signals:
			return
		}
	}
}

type Info struct {
	CompanyID uint   `json:"company_id"`
	ShopID    uint   `json:"shop_id"`
	ApiID     string `json:"api_id"`
	ApiSecret string `json:"api_secret"`
	FaceID    string `json:"face_id"`
	GroupID   string `json:"group_id"`
	PersonID  string `json:"person_id"`
}

func infoHandler(values []byte) {
	// step one: titan
	// step two: database event by person_id
	// step three : count and save into database

	var info interface{}
	if err := json.Unmarshal(values, info); err != nil {
		log.Println(err)
		return
	}

	var group *models.FrequentCustomerGroup
	var ok bool
	if ok, group = saveGroupInfo(info.(Info)); !ok {
		return
	}
	fetchDataByTitan(group, info.(Info))

}

// save group
func saveGroupInfo(info Info) (bool, *models.FrequentCustomerGroup) {
	var oneGroup models.FrequentCustomerGroup
	if dbError := database.POSTGRES.Where("company_id = ? AND shop_id = ?", info.CompanyID, info.ShopID).First(&oneGroup).Error; dbError != nil {
		oneGroup = models.FrequentCustomerGroup{
			CompanyID: info.CompanyID,
			ShopID:    info.ShopID,
			GroupUUID: utils.GenerateUUID(20),
		}
		if dbError := database.POSTGRES.Save(&oneGroup).Error; dbError != nil {
			return false, nil
		}
	}
	return true, &oneGroup
}

func fetchDataByTitan(group *models.FrequentCustomerGroup, info Info) {
	response, err := http.PostForm(titanParams.identificationURL, url.Values{
		"api_id":     {info.ApiID},
		"api_secret": {info.ApiSecret},
		"face_id":    {info.FaceID},
		"group_id":   {info.GroupID},
		"top":        {"20"},
	})

	if err != nil {
		return
	}
	defer response.Body.Close()

	var values interface{}
	responseByte, _ := ioutil.ReadAll(response.Body)
	if err := json.Unmarshal(responseByte, values); err != nil {
		return
	}
	var personIDs []string
	for _, value := range values.(titan.CandidateData).Candidates {
		personIDs = append(personIDs, value.PersonID)
	}
	personIDHandler(group, info, personIDs)

}

type result struct {
	PersonID string    `json:"person_id"`
	Day      time.Time `json:"day"`
}

type results []result

func personIDHandler(group *models.FrequentCustomerGroup, info Info, personIDs []string) {
	personIDString := strings.Join(personIDs, ",")
	now := time.Now()
	right := now.Format("2006-01-02 15:04:05")
	left := now.AddDate(0, -1, 0).Format("2006-01-02 15:04:05")
	sql := fmt.Sprintf(`SELECT person_id, date_trunc('day',max(capture_at)) as day FROM events WHERE person_id in (%s) AND capture_at BETWEEN %s AND %s ORDER BY capture_at desc`,
		personIDString, left, right)

	var resultsValues results

	database.POSTGRES.Raw(sql).Scan(&resultsValues)

	//personID
	{
		var onePerson models.FrequentCustomerPeople
		day, _ := time.Parse("2006-01-02 00:00:00", time.Now().Format("2006-01-02 00:00:00"))
		if dbError := database.POSTGRES.Where("person_id = ? AND date = ?", info.PersonID, day).First(&onePerson).Error; dbError != nil {
			onePerson = models.FrequentCustomerPeople{
				PersonID:                info.PersonID,
				FrequentCustomerGroupID: group.ID,
				Date:                    day,
				Interval:                uint(float64(time.Now().Sub(resultsValues[0].Day).Hours()/24) + 1),
				Frequency:               uint(len(resultsValues)),
			}
			database.POSTGRES.Save(&onePerson)
		}
	}
}
