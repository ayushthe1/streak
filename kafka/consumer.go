package kafka

import (
	"context"
	"encoding/json"
	"log"

	"github.com/IBM/sarama"
	"github.com/ayushthe1/streak/models"
	util "github.com/ayushthe1/streak/utils"
)

type consumer struct{}

func (consumer) Setup(sarama.ConsumerGroupSession) error   { return nil }
func (consumer) Cleanup(sarama.ConsumerGroupSession) error { return nil }
func (consumer) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		var notification models.Notification
		json.Unmarshal(message.Value, &notification)
		util.BroadcastKafkaEvent <- &notification
		sess.MarkMessage(message, "")
	}
	return nil
}

func ConsumeNotifications(ctx context.Context, brokers []string, groupID string, topics []string) {
	config := sarama.NewConfig()
	config.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategyRoundRobin()
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	consumerGroup, err := sarama.NewConsumerGroup(brokers, groupID, config)
	if err != nil {
		log.Fatalf("Failed to create consumer group: %s", err)
	}

	defer consumerGroup.Close()

	for {
		if err := consumerGroup.Consume(ctx, topics, &consumer{}); err != nil {
			log.Fatalf("Error from consumer: %s", err)
		}

		// Check if context is done
		if ctx.Err() != nil {
			log.Println("Context error: ", ctx.Err())
			return
		}
	}
}
