package kafka

import (
	"encoding/json"
	"log"

	"github.com/IBM/sarama"
)

var producer sarama.SyncProducer

func InitProducer(brokers []string) error {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 10
	config.Producer.Return.Successes = true

	var err error
	producer, err = sarama.NewSyncProducer(brokers, config)
	if err != nil {
		log.Fatal("not able to get producer :", err.Error())

	}

	return nil
}

// function to producer
func ProduceEventToKafka(topic string, message interface{}) error {
	msgBytes, err := json.Marshal(message)
	if err != nil {
		return err
	}

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(msgBytes),
	}

	_, _, err = producer.SendMessage(msg)
	if err != nil {
		log.Println("err 87")
		return err
	}
	log.Printf("Event produced to %s in ProduceNotification", topic)

	return nil
}
