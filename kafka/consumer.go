package kafka

import (
	"context"
	"encoding/json"
	"log"

	"github.com/IBM/sarama"
	"github.com/ayushthe1/streak/channels"
	"github.com/ayushthe1/streak/models"
)

var NotificationMsgType = "notification"
var ActivityMsgType = "activity"

var NotificationTopic = "notification_topic"
var ActivityTopic = "activity_topic"

type consumer struct{}

func (consumer) Setup(sarama.ConsumerGroupSession) error   { return nil }
func (consumer) Cleanup(sarama.ConsumerGroupSession) error { return nil }
func (consumer) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		log.Println("inside for loop in ConsumeClaim")

		// format of message.Value(in byte slice)
		// {
		// 	"type": "notification",
		// 	"data": {
		// 	  "id": 123,
		// 	  "message": "Hello, World!"
		// 	},
		// 	"timestamp": 1624892400
		//   }

		var msg map[string]interface{}
		err := json.Unmarshal(message.Value, &msg)
		if err != nil {
			log.Printf("Error unmarshalling")
			continue
		}

		// format of msg after unmarshalling
		// msg = map[string]interface{}{
		// 	"type": "notification",
		// 	"data": map[string]interface{}{
		// 	  "id": 123,
		// 	  "message": "Hello, World!",
		// 	},
		// 	"timestamp": 1624892400,
		//   }

		msgType, ok := msg["type"].(string)
		if !ok {
			log.Printf("msg type missing : %v", msgType)
		}

		switch msgType {

		case NotificationMsgType:
			var notification models.Notification
			err := json.Unmarshal(message.Value, &notification)
			if err != nil {
				log.Printf("Error unmarshalling activity: %v", err)
				continue
			}
			channels.BroadcastKafkaNotification <- &notification
			sess.MarkMessage(message, "")

		case ActivityMsgType:
			var activity models.ActivityEvent
			err := json.Unmarshal(message.Value, &activity)
			if err != nil {
				log.Printf("Error unmarshalling activity: %v", err)
				continue
			}
			channels.BroadcastKafkaActivity <- &activity
			sess.MarkMessage(message, "")

		default:
			log.Printf("Unhandled message type: %s", msgType)
		}

	}
	return nil
}

// function to consume kafka events from a given list of topics
func ConsumeMessages(ctx context.Context, brokers []string, groupID string, topics []string) {
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
			// After this the ConsumeClaim() function is called automatically
			log.Fatalf("Error from consumer: %s", err)
		}

		// Check if context is done
		if ctx.Err() != nil {
			log.Println("Context error: ", ctx.Err())
			return
		}
	}
}
