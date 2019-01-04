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
	"siren/pkg/database"
	"siren/pkg/titan"
	"time"

	"github.com/spf13/viper"

	cluster "github.com/bsm/sarama-cluster"
)

var headCountConsumerParams struct {
	brokers []string
	groupID string
	topics  []string
}

func consumerInit() {
	// todo: fix use config.fetchValue
	host := viper.Get(configs.ENV + ".kafka.host")
	port := viper.Get(configs.ENV + ".kafka.port")
	groupName := viper.GetString(configs.ENV + ".kafka.group")
	topic := viper.Get(configs.ENV + ".kafka.topic")
	log.Println(fmt.Sprintf("env: %s,host:%s, port: %s, groupID: %s, topic: %s", configs.ENV, host, port, groupName, topic))
	headCountConsumerParams.brokers = []string{fmt.Sprintf("%s:%s", host, port)}
	headCountConsumerParams.groupID = groupName
	headCountConsumerParams.topics = []string{topic.(string)}
}

func HeadCountConsumer() {
	consumerInit()
	StartForConsumer(
		headCountConsumerParams.brokers,
		headCountConsumerParams.groupID,
		headCountConsumerParams.topics,
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
	ApiID     string `json:"api_id"`
	ApiSecret string `json:"api_secret"`
	FaceID    string `json:"face_id"`
	GroupID   string `json:"group_id"`
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
	fetchDataByTitan("", info.(Info))

}

func fetchDataByTitan(link string, info Info) {
	response, err := http.PostForm(link, url.Values{
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
	var c chan string
	worker(c, values.(titan.CandidateData))
	personIDHandler(c)

}

func worker(c chan string, values titan.CandidateData) {
	for _, data := range values.Candidates {
		c <- data.PersonID
	}
}

type result struct {
	PersonID   string    `json:"person_id"`
	CapturedAt time.Time `json:"captured_at"`
}

type results []result

func personIDHandler(c chan string) {
	now := time.Now()
	right := now.Format("2006-01-02 15:04:05")
	left := now.AddDate(0, -1, 0).Format("2006-01-02 15:04:05")
	go func(left string, right string) {
		for {
			personID := <-c
			sql := fmt.Sprintf(`SELECT person_id, captured_at FROM events WHERE person_id = %s AND captured_at BETWEEN %s AND %s`,
				personID, left, right)
			var resultsValues results
			if dbError := database.POSTGRES.Raw(sql).Scan(&resultsValues).Error; dbError != nil {
				return
			}
		}
	}(left, right)
}
