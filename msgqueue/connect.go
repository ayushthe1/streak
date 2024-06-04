package msgqueue

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
)

var RabbitConn *amqp.Connection
var RabbitChannel *amqp.Channel

func init() {
	ConnectoRabbitMQ()
}

func ConnectoRabbitMQ() {
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

	log.Println("Connected to RabbitMQ")

}
