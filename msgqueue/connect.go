package msgqueue

import (
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	RabbitConn    *amqp.Connection
	RabbitChannel *amqp.Channel
	connLock      sync.Mutex
)

// func init() {
// 	ConnectoRabbitMQ()
// }

func ConnectoRabbitMQ() {
	log.Println("Inside ConnectToRabbitMQ")
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env files")
	}

	amqpServerURL := os.Getenv("AMQP_SERVER_URL")
	conn, err := amqp.Dial(amqpServerURL)

	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}

	RabbitConn = conn

	channel, err := RabbitConn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}

	RabbitChannel = channel

	log.Println("**********CONNECTED TO RABBITMQ**********")

	go func() {
		log.Println("Monitoring RabbitMQ connection")
		closeErr := <-RabbitChannel.NotifyClose(make(chan *amqp.Error))
		log.Printf("############# RabbitMQ CHANNEL CLOSED ###############: %v", closeErr)

		// Reconnect logic
		for {
			log.Println("Attempting to reconnect to RabbitMQ")
			// time.Sleep(5 * time.Second)
			connLock.Lock()
			ConnectoRabbitMQ()
			connLock.Unlock()
			log.Println("Reconnected to RabbitMQ")
			break
		}
	}()

}
