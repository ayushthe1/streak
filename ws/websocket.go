package ws

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/ayushthe1/streak/handler"
	"github.com/ayushthe1/streak/middleware"
	"github.com/ayushthe1/streak/models"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/websocket/v2"
)

type Client struct {
	Conn     *websocket.Conn
	Username string
}

type Message struct {
	Type string      `json:"type"`
	User string      `json:"user,omitempty"`
	Chat models.Chat `json:"chat,omitempty"`
}

// A map to keep track of active clients. The key is a pointer to a Client, and the value is a boolean
var clients = make(map[*Client]bool)

// A channel to broadcast chat messages to all connected clients.
var broadcast = make(chan *models.Chat)

// Handles incoming WebSocket requests.
func ServeWS(c *fiber.Ctx) error {
	log.Println("Inside ServeWS")
	// Ensure the connection is a WebSocket upgrade request
	if websocket.IsWebSocketUpgrade(c) {
		c.Locals("allowed", true)
		c.Locals("mycustomFlag", true)
		return websocket.New(func(conn *websocket.Conn) {
			client := &Client{
				Conn: conn,
			}
			clients[client] = true
			log.Printf("clients: %d, %v, %s", len(clients), clients, conn.RemoteAddr())

			defer func() {
				log.Printf("exiting: %s", conn.RemoteAddr().String())
				delete(clients, client)
				conn.Close()
			}()

			// Listen indefinitely for new messages coming through on the WebSocket connection
			receiver(client)
		})(c)
	}
	return fiber.ErrUpgradeRequired
}

func receiver(client *Client) {
	log.Println("Inside receiver")

	for {
		// read in a message
		// readMessage returns messageType, message, err
		// messageType: 1-> Text Message, 2 -> Binary Message
		_, p, err := client.Conn.ReadMessage()
		if err != nil {
			log.Println(err)
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
			client.Username = m.User
			fmt.Println("client successfully mapped", &client, client, client.Username)
		} else {
			fmt.Println("received message", m.Type, m.Chat)
			c := m.Chat //Copies the Chat field from the message to a local variable c.
			c.Timestamp = time.Now().Unix()

			// save in db
			id, err := handler.CreateChat(&c)
			if err != nil {
				log.Println("error while saving chat in DB", err)
				return
			}

			// is this really needed ?
			c.ID = id

			broadcast <- &c // Sends the chat message to the broadcast channel for broadcasting to other connected clients.
		}
	}
}

// broadcasting messages to all connected clients
// The function sends the message to all connected clients that are involved in the message
func broadcaster() {
	log.Println("Inside broadcaster")
	for {
		message := <-broadcast
		// send to every client that is currently connected
		fmt.Println("new message", message)

		for client := range clients {
			// send message only to involved users
			fmt.Println("username:", client.Username,
				"from:", message.From,
				"to:", message.To)

			if client.Username == message.From || client.Username == message.To {
				err := client.Conn.WriteJSON(message)
				if err != nil {
					log.Printf("Websocket error: %s", err)
					client.Conn.Close()
					delete(clients, client)
				}
			}
		}
	}
}

func Setup(app *fiber.App) {
	app.Use(middleware.IsAuthenticate)
	app.Get("/ws", ServeWS)
}

func StartWebSocketServer() {
	go broadcaster()

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
