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

// // function to consume messages from the rabbitMQ queue for a user. It starts for a user when he boots up
// func ConsumeMessages(ctx context.Context, client *Client) error {

// 	username := client.Username
// 	log.Println("going to consume")
// 	msgs, err := msgqueue.Consume(username)

// 	log.Println("This is the amqp091->delivery chan:", msgs)
// 	if err != nil {
// 		log.Println("Error inside consumeMessages while getting msg")
// 		return err
// 	}

// 	client.consumerChannel = msgs

// 	// Create a new context that can be canceled when the client disconnects
// 	consumeCtx, consumeCancel := context.WithCancel(ctx)
// 	client.consumeCancel = consumeCancel

// 	// Goroutine leak hapening
// 	go func() {
// 		defer func() {
// 			log.Println("Inside go func defers")
// 			if r := recover(); r != nil {
// 				log.Printf("Recovered in ConsumeMessages: %v", r)
// 			}
// 		}()

// 		for {
// 			select {
// 			// if this case dont execute in case of empty queue,then ConsumerTag is never sets
// 			case msg, ok := <-msgs:

// 				log.Println("executing case msg")
// 				if !ok {
// 					log.Println("Not ok")
// 					return
// 				}
// 				client.ConsumerTag = msg.ConsumerTag
// 				log.Println("consumer tah latest is : ", client.ConsumerTag)
// 				var chat models.Chat
// 				err := json.Unmarshal(msg.Body, &chat)
// 				if err != nil {
// 					log.Println("Failed to unmarshal chat message:", err)
// 					continue
// 				}
// 				log.Println("chat is : ", chat)

// 				// Push the message to the WebSocket connection
// 				log.Println("Going to write to conn")
// 				err = client.Conn.WriteJSON(chat)
// 				if err != nil {
// 					log.Println("Failed to write message to WebSocket inside ConsumeMessage:", err)
// 					return
// 				}
// 				msg.Ack(false)
// 			case <-consumeCtx.Done():
// 				log.Println("ConsumeMessages goroutine canceled")
// 				return
// 			}

// 		}
// 	}()

// 	log.Println("At end of ConsumeMessages")
// 	return nil

// }

// func ServeWS(c *fiber.Ctx) error {
// 	log.Println("Inside ServeWS")
// 	// Ensure the connection is a WebSocket upgrade request
// 	if websocket.IsWebSocketUpgrade(c) {
// 		c.Locals("allowed", true)
// 		c.Locals("mycustomFlag", true)
// 		return websocket.New(func(conn *websocket.Conn) {

// 			// Create a context for canceling the goroutine
// 			ctx, cancel := context.WithCancel(context.Background())
// 			client := &Client{
// 				Conn:   conn,
// 				cancel: cancel,
// 			}

// 			clients.Lock()
// 			clients.m[client] = true
// 			clients.Unlock()

// 			log.Printf("clients: %d, %v, %s", len(clients.m), clients.m, conn.RemoteAddr())

// 			defer func() {
// 				log.Printf("Closing the connection: %s", conn.RemoteAddr().String())

// 				if err != nil {
// 					log.Printf("error deleting user %s from redis", client.Username)
// 				}
// 				cancel() // Cancel the goroutine
// 				client.consumeCancel()
// 				if client.consumerChannel != nil {
// 					client.Conn.Close()
// 					log.Println("consumer tag before vlosing : ", client.ConsumerTag)
// 					msgqueue.RabbitChannel.Cancel(client.ConsumerTag, true)
// 					log.Printf("Consumer tag for user %s is %s ", client.Username, client.ConsumerTag)
// 					log.Println("Delivery stopped")
// 				}
// 			}()

// 			// Listen indefinitely for new messages coming through on the WebSocket connection
// 			receiveMessages(ctx, client)
// 		})(c)
// 	}
// 	return fiber.ErrUpgradeRequired
// }

// // Listen for all messages coming on the websocket connection
// func receiveMessages(ctx context.Context, client *Client) {
// 	log.Println("Inside receiver")

// 	for {
// 		log.Println("For loop in receiveMsg")
// 		_, p, err := client.Conn.ReadMessage()
// 		if err != nil {
// 			client.cancel() // Cancel the context to stop the goroutine
// 			log.Println("HII")
// 			log.Println("some error", err.Error())
// 			return
// 		}

// 		m := &Message{}

// 		err = json.Unmarshal(p, m)
// 		if err != nil {
// 			log.Println("error while unmarshaling chat", err)
// 			continue
// 		}

// 		fmt.Println("host", client.Conn.RemoteAddr())
// 		//  If the type is "bootup", it sets the client's username to the user in the message.
// 		if m.Type == "bootup" {
// 			// do mapping on bootup
// 			log.Println("m.user is : ", m.User)
// 			client.Username = m.User
// 			var testChat models.Chat
// 			testChat.Msg = "TEST"
// 			sendMessageToQueue(client.Username, &testChat)

// 			// Start consuming messages from RabbitMQ for the user
// 			go func() {
// 				err := ConsumeMessages(ctx, client)
// 				if err != nil {
// 					log.Println("Error while consuming messages:", err)
// 					os.Exit(1)
// 				}
// 			}()

// 		} else if m.Type == "chat" {
// 			// deliver message to the receiver
// 			c := m.Chat
// 			c.Timestamp = time.Now().Unix()
// 			deliverMessageToUser(recieverUsername, &c)

// 		} else {
// 			log.Printf("Invalid Type : %s", m.Type)
// 			client.Conn.WriteJSON(fmt.Sprintf("error : invalid msg type : %s", m.Type))
// 		}
// 	}

// }
