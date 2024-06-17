package kafka

import (
	"context"
	"encoding/json"
	"log"

	"github.com/IBM/sarama"
	"github.com/ayushthe1/streak/channels"
	"github.com/ayushthe1/streak/handler"
	"github.com/ayushthe1/streak/models"
)

var NotificationMsgType = "notification" // sent only to a specific user
var ActivityMsgType = "activity"         // sent to all users
var ChatMsgType = "chat"                 // chat is inserted in DB

var NotificationTopic = "notification_topic" // all notification events(event meant for any specific user) published to this topic
var ActivityTopic = "activity_topic"         // all activity events(events that will be sent to all use3rs) published to this topic
var ChatTopic = "chat_topic"                 // all chat messages will be published to this topic

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
				log.Printf("Error unmarshalling activity: %v", err) // Unmarshal(nil *models.Chat)
				continue
			}
			channels.BroadcastKafkaActivity <- &activity
			sess.MarkMessage(message, "")

		case ChatMsgType: //TODO: maybe use channel later to avoid blocking ?
			var chatevent models.ChatEvent
			err := json.Unmarshal(message.Value, &chatevent)
			log.Println("Chat msg received in consumer is : ", chatevent)
			if err != nil {
				log.Printf("error unmarsahlling the chat msg in kafka %v ", err)
				continue
			}

			// save the chat in DB
			chat := chatevent.ChatMsg
			_, err = handler.CreateChat(chat)
			if err != nil {
				log.Fatalf("unabe to save chat msg in db : %s", err.Error())
			}
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
