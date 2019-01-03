package kafka

import (
	"encoding/json"
	"fmt"
	"log"
	"siren/configs"
	"strings"
	"time"

	"github.com/spf13/viper"

	"github.com/Shopify/sarama"
)

var headCountProducerParams struct {
	brokers string
}

func producerInit() {
	host := viper.Get(configs.ENV + ".kafka.host")
	port := viper.Get(configs.ENV + ".kafka.port")
	log.Println(fmt.Sprintf("host:%s, prot:%s", host, port))
	headCountProducerParams.brokers = fmt.Sprintf("%s:%s", host, port)
}
func HeadCountProducer() Server {
	producerInit()
	return Start(headCountProducerParams.brokers)
}

var ProducerServer Server

type Server struct {
	AccessLogProducer sarama.AsyncProducer
}

func Start(brokerString string) Server {
	brokerList := strings.Split(brokerString, ",")
	server := Server{
		AccessLogProducer: newAccessLogProducer(brokerList),
	}
	ProducerServer = server
	return ProducerServer

}

func (s *Server) Close() error {
	if err := s.AccessLogProducer.Close(); err != nil {
		log.Println("Failed to shut down access log producer cleanly", err)
	}
	return nil
}

type accessLogEntry struct {
	encoded []byte `json:"encoded"`
	err     error  `json:"err"`
}

func (a *accessLogEntry) ensureEncoded() {
	if a.encoded == nil && a.err == nil {
		a.encoded, a.err = json.Marshal(a)
	}
}

func (a *accessLogEntry) Length() int {
	a.ensureEncoded()
	return len(a.encoded)
}

func (a *accessLogEntry) Encode() ([]byte, error) {
	a.ensureEncoded()
	return a.encoded, a.err
}

func (s *Server) WithAccessLog(topic string, key string, value string) {
	s.AccessLogProducer.Input() <- &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(key),
		Value: &accessLogEntry{encoded: []byte(value)},
	}
}

func newAccessLogProducer(brokerList []string) sarama.AsyncProducer {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForLocal       // Only wait for the leader to ack
	config.Producer.Compression = sarama.CompressionSnappy   // Compress messages
	config.Producer.Flush.Frequency = 500 * time.Millisecond // Flush batches every 500ms
	producer, err := sarama.NewAsyncProducer(brokerList, config)
	if err != nil {
		log.Fatalln("Failed to start Sarama producer:", err)
	}

	// We will just log to STDOUT if we're not able to produce messages.
	// Note: messages will only be returned here after all retry attempts are exhausted.
	go func() {
		for err := range producer.Errors() {
			log.Println("Failed to write access log entry:", err)
		}
	}()

	return producer
}
