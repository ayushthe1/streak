package ws

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/ayushthe1/streak/handler"
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

// A map to keep track of active clients. The key is a pointer to a Client, and the value is a boolean
var clients = struct {
	sync.Mutex
	m map[*Client]bool
}{m: make(map[*Client]bool)}

// var clients = make(map[string]*Client)

// A channel to broadcast chat messages to all connected clients.
var broadcast = make(chan *models.Chat)

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

				clients.Lock()
				delete(clients.m, client)
				clients.Unlock()

				presence.setUserOffline(client.Username)
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
		// read in a message
		// readMessage returns messageType, message, err
		// messageType: 1-> Text Message, 2 -> Binary Message
		log.Println("For loop in receiveMsg")
		_, p, err := client.Conn.ReadMessage()
		if err != nil {
			client.cancel() // Cancel the context to stop the goroutine
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
			presence.setUserOnline(client.Username)
			fmt.Println("client successfully mapped", &client, client, client.Username)

			// deliver messages to user already stored in queue
			log.Println("Calling ConsumeMessags function by :", client.Username)

			// push a test msg to the queue
			var testChat models.Chat
			testChat.Msg = "TEST"
			sendMessageToQueue(client.Username, &testChat)

			err := ConsumeMessages(ctx, client)
			if err != nil {
				log.Println("Error while consuming messages:", err)
				return
			}
		} else {
			// deliver message to the receiver

			log.Println("m.User : ", m.User)
			fmt.Println("received message : ", m.Type, m.Chat)
			c := m.Chat //Copies the Chat field from the message to a local variable c.
			c.Timestamp = time.Now().Unix()

			// save in db
			_, err := handler.CreateChat(&c)
			if err != nil {
				log.Println("error while saving chat in DB", err)
				client.cancel() // Cancel the context to stop the goroutine
				return
			}

			// // is this really needed ?
			// log.Print("This is ID of C after chat is created", c.Id)
			// log.Println("This is the id value returned :", id)
			// // c.ID = id

			recieverUsername := c.To

			if presence.isUserOnline((recieverUsername)) {
				broadcast <- &c // Sends the chat message to the broadcast channel for broadcasting to other client (set in chat.To field only when he is online).
			} else {
				log.Printf("User %s is not online currently. Sending message to the queue", recieverUsername)

				sendMessageToQueue(recieverUsername, &c)
			}

		}
	}
	log.Println("At end of receive msg")
}

// function to deliver message to the user when he is online
func deliverMessageWhenUserIsOnline() {
	log.Println("Inside broadcaster")
	for {
		log.Println("For in when user is online")
		message := <-broadcast
		// send to every client that is currently connected
		fmt.Println("new message", message)

		clients.Lock()
		for client := range clients.m {
			// send message only to involved users
			fmt.Println("username:", client.Username,
				"from:", message.From,
				"to:", message.To)

			if client.Username == message.From || client.Username == message.To {
				err := client.Conn.WriteJSON(message) // send the msg
				if err != nil {
					log.Printf("Websocket error: %s", err)
					client.Conn.Close()
					delete(clients.m, client)
				}
			}
		}
		clients.Unlock()
	}
}

// function to consume messages from the rabbitMQ queue for a user
func ConsumeMessages(ctx context.Context, client *Client) error {

	username := client.Username
	log.Println("going to consume")
	msgs, err := msgqueue.Consume(username)

	log.Println("This is the amqp091->delivery chan:", msgs)
	if err != nil {
		log.Println("Error inside consumeMessages while getting msg")
		return err
	}

	client.consumerChannel = msgs

	// Create a new context that can be canceled when the client disconnects
	consumeCtx, consumeCancel := context.WithCancel(ctx)
	client.consumeCancel = consumeCancel

	// Goroutine leak hapening
	go func() {
		defer func() {
			log.Println("Inside go func defers")
			if r := recover(); r != nil {
				log.Printf("Recovered in ConsumeMessages: %v", r)
			}
		}()

		log.Println("len of msg :", len(msgs))
		for {
			log.Println("executing for")
			select {
			// if this case dont execute in case of empty queue,then ConsumerTag is never sets
			case msg, ok := <-msgs:

				log.Println("executing case msg")
				if !ok {
					log.Println("Not ok")
					return
				}
				client.ConsumerTag = msg.ConsumerTag
				log.Println("consumer tah latest is : ", client.ConsumerTag)
				var chat models.Chat
				err := json.Unmarshal(msg.Body, &chat)
				if err != nil {
					log.Println("Failed to unmarshal chat message:", err)
					continue
				}
				log.Println("chat is : ", chat)

				// Push the message to the WebSocket connection
				log.Println("Going to write to conn")
				err = client.Conn.WriteJSON(chat)
				if err != nil {

					log.Println("Failed to write message to WebSocket inside ConsumeMessage:", err)
					return
				}
				msg.Ack(false)
			case <-consumeCtx.Done():
				log.Println("ConsumeMessages goroutine canceled")
				return
			}

		}
	}()

	log.Println("At end of ConsumeMessages")
	return nil

}

// function to send message to the queue if user is offline
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
	go deliverMessageWhenUserIsOnline()

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