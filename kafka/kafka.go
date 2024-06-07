package kafka

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/net/context"
)

func StartKafkaConsumer() {
	brokers := []string{"localhost:9092"}
	if err := InitProducer(brokers); err != nil {
		log.Fatalf("Failed to initialize Kafka producer: %s", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go ConsumeNotifications(ctx, brokers, "notification_group", []string{"notification_topic"})

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop // Wait for a signal

	log.Println("Shutting down gracefully...")
	// cancel() called

}
