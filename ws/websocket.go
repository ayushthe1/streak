package ws

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/ayushthe1/streak/channels"
	"github.com/ayushthe1/streak/online"

	"github.com/ayushthe1/streak/kafka"
	"github.com/ayushthe1/streak/models"
	"github.com/ayushthe1/streak/msgqueue"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/websocket/v2"
	"github.com/rabbitmq/amqp091-go"
)

type Client struct {
	Conn            *websocket.Conn
	Username        string
	cancel          context.CancelFunc      // for cancelling the goroutine
	consumeCancel   context.CancelFunc      // for canceling the consume goroutine
	consumerChannel <-chan amqp091.Delivery // Channel for consuming messages
	ConsumerTag     string
}

type Message struct {
	Type string      `json:"type"`
	User string      `json:"user,omitempty"`
	Chat models.Chat `json:"chat,omitempty"`
}

// A map to keep track of active clients. The key is a pointer to a Client, and the value is a boolean. We don't really need it as now we're using Redis for storing online users
var clients = struct {
	sync.Mutex
	m map[*Client]bool
}{m: make(map[*Client]bool)}

// var clients = make(map[string]*Client)

var presence PresenceService

func init() {
	onlineUsers := map[string]bool{}
	mu := &sync.Mutex{}
	presence = PresenceService{
		onlineUsers: onlineUsers,
		mu:          mu,
	}
}

// Handles incoming WebSocket requests.
func ServeWS(c *fiber.Ctx) error {
	log.Println("Inside ServeWS")
	// Ensure the connection is a WebSocket upgrade request
	if websocket.IsWebSocketUpgrade(c) {
		c.Locals("allowed", true)
		c.Locals("mycustomFlag", true)
		return websocket.New(func(conn *websocket.Conn) {

			// Create a context for canceling the goroutine
			ctx, cancel := context.WithCancel(context.Background())
			client := &Client{
				Conn:   conn,
				cancel: cancel,
			}

			clients.Lock()
			clients.m[client] = true
			clients.Unlock()

			log.Printf("clients: %d, %v, %s", len(clients.m), clients.m, conn.RemoteAddr())

			defer func() {
				log.Printf("Closing the connection: %s", conn.RemoteAddr().String())

				//TODO: not needed, remove later
				presence.setUserOffline(client.Username)

				err := online.DeleteUserFromRedis(client.Username)
				if err != nil {
					log.Printf("error deleting user %s from redis", client.Username)
				}
				cancel() // Cancel the goroutine
				client.consumeCancel()
				if client.consumerChannel != nil {
					client.Conn.Close()
					log.Println("consumer tag before vlosing : ", client.ConsumerTag)
					msgqueue.RabbitChannel.Cancel(client.ConsumerTag, true)
					log.Printf("Consumer tag for user %s is %s ", client.Username, client.ConsumerTag)
					log.Println("Delivery stopped")
				}
				// conn.Close()
				logoutActivity := models.ActivityEvent{
					Type:      kafka.ActivityMsgType,
					Username:  client.Username,
					Action:    "log out",
					Timestamp: time.Now().Unix(),
					Details:   fmt.Sprintf("User %s just logged out", client.Username),
				}

				if err := kafka.ProduceEventToKafka(kafka.ActivityTopic, logoutActivity); err != nil {
					log.Printf("error producing logout activity %v", err)

				}

				clients.Lock()
				delete(clients.m, client)
				clients.Unlock()

				log.Println("At end of defer close")
			}()

			// Listen indefinitely for new messages coming through on the WebSocket connection
			receiveMessages(ctx, client)
		})(c)
	}
	return fiber.ErrUpgradeRequired
}

// Listen for all messages coming on the websocket connection
func receiveMessages(ctx context.Context, client *Client) {
	log.Println("Inside receiver")

	for {
		log.Println("For loop in receiveMsg")
		_, p, err := client.Conn.ReadMessage()
		if err != nil {
			// client.cancel() // Cancel the context to stop the goroutine
			log.Println("HII")
			log.Println("some error", err.Error())
			return
		}

		m := &Message{}

		err = json.Unmarshal(p, m)
		if err != nil {
			log.Println("error while unmarshaling chat", err)
			continue
		}

		fmt.Println("host", client.Conn.RemoteAddr())
		//  If the type is "bootup", it sets the client's username to the user in the message.
		if m.Type == "bootup" {
			// do mapping on bootup
			log.Println("m.user is : ", m.User)
			client.Username = m.User

			// TODO: not needed, remove later
			presence.setUserOnline(client.Username)

			err := online.AddUserToRedis(client.Username)
			if err != nil {
				log.Println("error adding user to redis")
			}

			// Produce welcome message to Kafka
			log.Println("Producing welcome msg to kafka")
			welcomeEvent := models.Notification{
				Type:     kafka.NotificationMsgType,
				Username: client.Username,
				Message:  "Welcome to the chat!",
			}
			if err := kafka.ProduceEventToKafka(kafka.NotificationTopic, welcomeEvent); err != nil {
				log.Printf("Failed to produce welcome message: %s", err)
			}

			// publish login activiy to kafka which will be sent to all users
			loginActivity := models.ActivityEvent{
				Type:      kafka.ActivityMsgType,
				Username:  client.Username,
				Action:    "login",
				Timestamp: time.Now().Unix(),
				Details:   fmt.Sprintf("User %s just logged in", client.Username),
			}
			if err := kafka.ProduceEventToKafka(kafka.ActivityTopic, loginActivity); err != nil {
				log.Printf("Failed to produce loginActivity : %v", loginActivity)
			}

			// // deliver messages to user already stored in queue
			// log.Println("Calling ConsumeMessags function by :", client.Username)

			// push a test msg to the queue
			var testChat models.Chat
			testChat.Msg = "TEST"
			sendMessageToQueue(client.Username, &testChat)

			// Start consuming messages from RabbitMQ for the user
			go func() {
				err := ConsumeMessages(ctx, client)
				if err != nil {
					log.Println("Error while consuming messages:", err)
					os.Exit(1)
				}
			}()

			go deliverNotificationWhenUserIsOnline(ctx, client)

		} else if m.Type == "chat" {
			// deliver message to the receiver
			c := m.Chat
			c.Timestamp = time.Now().Unix()
			log.Println("The chat msg received is :", c.Msg)

			chatMsg := c
			// Publish the chat message to kafka
			chatEvent := models.ChatEvent{
				Type:    kafka.ChatMsgType,
				ChatMsg: &chatMsg,
			}
			if err := kafka.ProduceEventToKafka(kafka.ChatTopic, chatEvent); err != nil {
				log.Printf("Failed to produce chatevent : %v", err)
			}

			recieverUsername := c.To

			// Irrespective of whether user is online or not, send the message to RabbitMQ
			deliverMessageToUser(recieverUsername, &c)

		} else {
			log.Printf("Invalid Type : %s", m.Type)
			client.Conn.WriteJSON(fmt.Sprintf("error : invalid msg type : %s", m.Type))
		}
	}

}

func deliverMessageToUser(toUsername string, chat *models.Chat) {

	// publish the message to rabbimq with the to_useranme as queue name
	err := sendMessageToQueue(toUsername, chat)
	if err != nil {
		log.Fatalf("unable to send message to queue : %v", err.Error())
	}

	// if user is online ,also send them a notificdation
	if online.IsUserOnline(toUsername) {
		notificationMsg := models.Notification{
			Type:     kafka.NotificationMsgType,
			Username: toUsername,
			Message:  "You have a new message from " + chat.From,
		}
		if err := kafka.ProduceEventToKafka(kafka.NotificationTopic, notificationMsg); err != nil {
			log.Printf("Failed to produce notification message: %s", err)
		}
	}

	// publish event to kafka about chat message which will be consumed by all users
	chatActivity := models.ActivityEvent{
		Type:      kafka.ActivityMsgType,
		Username:  chat.From,
		Action:    "sent a message",
		Timestamp: time.Now().Unix(),
		Details:   fmt.Sprintf("user %s just sent a msg to user %s", chat.From, chat.To),
	}

	if err := kafka.ProduceEventToKafka(kafka.ActivityTopic, chatActivity); err != nil {
		log.Printf("error producing chat activity %v", err)

	}

}

// Consumer function to take messages from kafka and sent it to only a specific user for whom the notification is meant for
func deliverNotificationWhenUserIsOnline(ctx context.Context, client *Client) {
	log.Println("Inside deliver notification")

	for {
		select {
		case <-ctx.Done():
			log.Println("Context canceled, stopping notification delivery for client:", client.Username)
			return
		case notification := <-channels.BroadcastKafkaNotification:
			log.Println("New notification to send:", notification)
			if client.Username == notification.Username {
				err := client.Conn.WriteJSON(notification)
				if err != nil {
					log.Printf("WebSocket error: %s", err)
					client.Conn.Close()
					online.DeleteUserFromRedis(client.Username)
					return
				}
			}
		}
	}
}

// TODO: change the logic here ,so it works for scalable application. It currently looks for clients from in-memory map which wont work when application scales
// Consumer function to take activity messages from kafka and send it to all online users
func deliverActivityToOnlineUsers() {
	for {
		activity := <-channels.BroadcastKafkaActivity
		log.Println("new activity to send :", activity)

		clients.Lock()
		for client := range clients.m {

			err := client.Conn.WriteJSON(activity) // send the msg
			if err != nil {
				log.Printf("Websocket error: %s", err)
				client.Conn.Close()
				delete(clients.m, client)
			}
		}
		clients.Unlock()
	}

}

// function to consume messages from the rabbitMQ queue for a user. It starts for a user when he boots up
func ConsumeMessages(ctx context.Context, client *Client) error {
	username := client.Username
	log.Println("Going to consume")
	msgs, err := msgqueue.Consume(username)

	if err != nil {
		log.Println("Error inside ConsumeMessages while getting msg:", err)
		return err
	}

	client.consumerChannel = msgs

	// Create a new context that can be canceled when the client disconnects
	consumeCtx, consumeCancel := context.WithCancel(ctx)
	client.consumeCancel = consumeCancel
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			select {
			case msg, ok := <-msgs:
				if !ok {
					log.Println("Message channel closed")
					return
				}

				client.ConsumerTag = msg.ConsumerTag
				log.Println("Consumer tag latest is:", client.ConsumerTag)
				var chat models.Chat
				err := json.Unmarshal(msg.Body, &chat)
				if err != nil {
					log.Println("Failed to unmarshal chat message:", err)
					continue
				}
				log.Println("Chat is:", chat)

				// Push the message to the WebSocket connection
				log.Println("Going to write to conn")
				err = client.Conn.WriteJSON(chat)
				if err != nil {
					log.Println("Failed to write message to WebSocket inside ConsumeMessage:", err)
					return
				}
				msg.Ack(false) // Acknowledge only this message ,True means acknowledge this and all the previous message.
			case <-consumeCtx.Done():
				log.Println("ConsumeMessages goroutine canceled")
				return
			}
		}
	}()

	wg.Wait()
	log.Println("At end of ConsumeMessages")
	return nil
}

// function to publish message to the queue.
// username is the receivers username and is the name of queue to which the msg will be published
func sendMessageToQueue(username string, chat *models.Chat) error {
	err := msgqueue.Publish(username, chat)

	if err != nil {
		log.Printf("Error while sending msg to queue")
		return err
	}

	log.Printf("At end of sendMessageToQueue. msg sent successfully")
	return nil
}

func Setup(app *fiber.App) {
	// app.Use(middleware.IsAuthenticate)
	app.Get("/ws", ServeWS)
}

func StartWebSocketServer() {
	go deliverActivityToOnlineUsers()

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:4000 ",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: true,
	}))

	Setup(app)
	port := "3001"
	log.Println("Web Socket Server listening on port", port)
	app.Listen(":" + port)

}
