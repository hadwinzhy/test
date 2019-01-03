package kafka

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"siren/configs"

	"github.com/spf13/viper"

	cluster "github.com/bsm/sarama-cluster"
)

var headCountCustomerParams struct {
	brokers []string
	groupID string
	topics  []string
}

func customerInit() {
	host := viper.Get(configs.ENV + ".kafka.host")
	port := viper.Get(configs.ENV + ".kafka.port")
	groupName := viper.GetString(configs.ENV + ".kafka.group")
	topic := viper.Get(configs.ENV + ".kafka.topic")
	log.Println(fmt.Sprintf("env: %s,host:%s, port: %s, groupID: %s, topic: %s", configs.ENV, host, port, groupName, topic))
	headCountCustomerParams.brokers = []string{fmt.Sprintf("%s:%s", host, port)}
	headCountCustomerParams.groupID = groupName
	headCountCustomerParams.topics = []string{topic.(string)}
}

func HeadCountCustomer() {
	customerInit()
	StartForCustomer(
		headCountCustomerParams.brokers,
		headCountCustomerParams.groupID,
		headCountCustomerParams.topics,
	)
}

func StartForCustomer(brokers []string, groupID string, topics []string) {
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

func infoHandler(values []byte) {
	// step one: titan
	// step two: database event by person_id
	// step three : count and save into database
	fmt.Println(string(values))
}
