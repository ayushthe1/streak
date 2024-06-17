package msgqueue

import (
	"encoding/json"
	"log"

	"github.com/ayushthe1/streak/models"
	amqp "github.com/rabbitmq/amqp091-go"
)

// function to publish message to RabbitMQ. The queueName should be the username of the user
func Publish(queueName string, chat *models.Chat) error {
	// ConnectoRabbitMQ()
	connLock.Lock()
	defer connLock.Unlock()

	_, err := RabbitChannel.QueueDeclare(
		queueName, // name of the queue
		true,
		false,
		false,
		false,
		nil, //  arguments
	)

	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
		return err
	}

	message, err := json.Marshal(chat)
	if err != nil {
		log.Println("Failed to marshal chat")
		return err
	}

	err = RabbitChannel.Publish(
		"",
		queueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        message,
		},
	)

	if err != nil {
		log.Println("Error publishing to channel")
		return err
	}

	log.Printf("Published message successfully : %v", chat)
	return nil

}

func Consume(queueName string) (<-chan amqp.Delivery, error) {

	connLock.Lock()
	defer connLock.Unlock()

	// log.Println("Inside Consume func : Is the channel closed ?", RabbitChannel.IsClosed())
	// log.Println("Inside Consume func : Is the conn closed ?", RabbitConn.IsClosed())

	log.Printf("Inside Consume func, going to consume now for user %s ", queueName)
	_, err := RabbitChannel.QueueDeclare(
		queueName, // name of the queue
		true,
		false,
		false,
		false,
		nil, //  arguments
	)

	if err != nil {
		log.Println("err :", err.Error())
		return nil, err
	}

	return RabbitChannel.Consume(
		queueName,
		"",
		false, // auto-ack -  Setting this to false means the consumer must explicitly acknowledge each message after processing it.
		false,
		false,
		false,
		nil,
	)

}
