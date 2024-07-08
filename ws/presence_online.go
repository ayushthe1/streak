package ws

import (
	"log"
	"sync"
)

type PresenceService struct {
	onlineUsers map[string]bool
	mu          *sync.Mutex
}

func (p *PresenceService) setUserOnline(username string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.onlineUsers[username] = true

}

func (p *PresenceService) setUserOffline(username string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	log.Printf("Going to set user %s offline", username)

	if p.onlineUsers == nil {
		log.Println("onlineUsers map is nil")
		return
	}

	_, ok := p.onlineUsers[username]
	if !ok {
		log.Printf("User %s is already offline", username)
		return
	}
	delete(p.onlineUsers, username)
}

func (p *PresenceService) isUserOnline(Username string) bool {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.onlineUsers == nil {
		log.Println("onlineUsers map in iUO is nil")
		return false
	}
	return p.onlineUsers[Username]

}

// package ws

// import (
// 	"sync"

// 	"github.com/gofiber/websocket/v2"
// )

// type UserPresence struct {
// 	username string
// 	conn     *websocket.Conn
// 	active   bool
// }

// type PresenceService struct {
// 	onlineUsers map[string]UserPresence
// 	mu          *sync.Mutex
// }

// func (p *PresenceService) setUserOnline(username string, conn *websocket.Conn) {
// 	p.mu.Lock()
// 	defer p.mu.Unlock()
// 	p.onlineUsers[username] = UserPresence{
// 		username: username,
// 		conn:     conn,
// 	}

// }

// func (p *PresenceService) setUserOffline(username string) {
// 	p.mu.Lock()
// 	defer p.mu.Unlock()
// 	delete(p.onlineUsers, username)
// }

// func (p *PresenceService) isUserOnline(Username string) bool {
// 	p.mu.Lock()
// 	defer p.mu.Unlock()
// 	return p.onlineUsers[Username].active

// }
