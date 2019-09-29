package kafka

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"siren/configs"
	"siren/pkg/logger"
	"time"

	"siren/venus/venus-model/models"

	"siren/src/workers"

	"github.com/bsm/sarama-cluster"
)

type CountFrequentConsumerParamsType struct {
	brokers []string
	groupID string
	topics  []string
	handler func([]byte, []byte)
}

var MallCountFrequentConsumerParams CountFrequentConsumerParamsType
var StoreCountFrequentConsumerParams CountFrequentConsumerParamsType

var titanParams struct {
	identificationURL string
	groupCreateURL    string
	groupAddPerson    string
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
	titanParams.groupAddPerson = fmt.Sprintf(configs.FetchFieldValue("TitanHOST") + "/groups/add_person")
	// todo: rule number
	ruleNumber.high = 3
	ruleNumber.low = 2
}

func CountFrequentConsumer() {
	consumerInit()
	go MallCountFrequentConsumerParams.StartConsumer()
	StoreCountFrequentConsumerParams.StartConsumer()
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
	workerHash := time.Now().UnixNano()
	for {
		select {
		case msg, ok := <-consumer.Messages():
			if ok {
				fmt.Fprintf(os.Stdout, "%d %s/%d/%d\t%s\t%s\n", workerHash, msg.Topic, msg.Partition, msg.Offset, msg.Key, msg.Value)
				params.handler(msg.Key, msg.Value)
				consumer.MarkOffset(msg, "") // mark message as processed
			}
		case <-signals:
			return
		}
	}
}

type InfoForKafkaProducer struct {
	CompanyID   uint   `json:"company_id"`
	FaceID      string `json:"face_id"`
	PersonID    string `json:"person_id"`
	CapturedAt  int64  `json:"captured_at"`
	EventID     uint   `json:"event_id"`
	EventStatus string `json:"event_status"`
}

func mallInfoHandler(key []byte, values []byte) {
	var info InfoForKafkaProducer
	if err := json.Unmarshal(values, &info); err != nil {
		log.Println(err)
		return
	}

	var group *models.FrequentCustomerGroup
	var ok bool
	nowSaveGroupInfo := time.Now()
	logger.Info("statistic time", "save group info", "start")
	if ok, group = saveGroupInfo(info.CompanyID); !ok {
		return
	}
	logger.Info("statistic time", "save group info", "count time", time.Now().Sub(nowSaveGroupInfo).Nanoseconds()/1000000)

	nowGroupAddPerson := time.Now()
	logger.Info("statistic time", "titan group add person", "start")
	// todo: personID or faceID?
	if info.PersonID != "" {
		go titanGroupAddPerson(group.GroupUUID, info.PersonID)
	}
	logger.Info("statistic time", "titan group add person", "count time", time.Now().Sub(nowGroupAddPerson).Nanoseconds()/1000000)
	nowFetchDataByTitan := time.Now()
	logger.Info("statistic time", "fetch data by titan", "start")
	fetchDataByTitan(group, info)
	logger.Info("statistic time", "fetch data by titan", "count time", time.Now().Sub(nowFetchDataByTitan).Nanoseconds()/1000000)

}

type StoreInfo struct {
	CompanyID uint   `json:"company_id"`
	ShopID    uint   `json:"shop_id"`
	PersonID  string `json:"person_id"`
	CaptureAt int64  `json:"capture_at"`
	EventID   uint   `json:"event_id"`
}

func storeInfoHandler(key []byte, values []byte) {
	var info StoreInfo
	if err := json.Unmarshal(values, &info); err != nil {
		log.Println(err)
		return
	}

	switch string(key) {
	case "remove":
		workers.RemoveFrequentCustomerHandler(info.PersonID)
	default:
		workers.StoreFrequentCustomerHandler(info.CompanyID, info.ShopID, info.PersonID, info.CaptureAt, info.EventID)
	}

}
