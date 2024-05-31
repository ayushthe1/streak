package ws

import (
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
	delete(p.onlineUsers, username)
}

func (p *PresenceService) isUserOnline(Username string) bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.onlineUsers[Username]

}
