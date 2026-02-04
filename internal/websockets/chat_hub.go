package websockets

import (
	"sync"

	"github.com/gorilla/websocket"
)

// ChatHub manages all active chat connections
type ChatHub struct {
	// Map UserID -> List of Connections (One user might have 2 tabs open)
	Clients map[uint]map[*websocket.Conn]bool

	// Actions
	Register   chan *ChatClient
	Unregister chan *ChatClient
	Broadcast  chan *ChatMessage

	mu sync.RWMutex
}

type ChatClient struct {
	UserID uint
	Conn   *websocket.Conn
}

type ChatMessage struct {
	ToUserID   uint
	FromUserID uint
	Content    string
	Time       string
}

func NewChatHub() *ChatHub {
	return &ChatHub{
		Clients:    make(map[uint]map[*websocket.Conn]bool),
		Register:   make(chan *ChatClient),
		Unregister: make(chan *ChatClient),
		Broadcast:  make(chan *ChatMessage),
	}
}

func (h *ChatHub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.mu.Lock()
			if _, ok := h.Clients[client.UserID]; !ok {
				h.Clients[client.UserID] = make(map[*websocket.Conn]bool)
			}
			h.Clients[client.UserID][client.Conn] = true
			h.mu.Unlock()

		case client := <-h.Unregister:
			h.mu.Lock()
			if _, ok := h.Clients[client.UserID]; ok {
				delete(h.Clients[client.UserID], client.Conn)
				client.Conn.Close()
				// Clean up empty maps
				if len(h.Clients[client.UserID]) == 0 {
					delete(h.Clients, client.UserID)
				}
			}
			h.mu.Unlock()

		case msg := <-h.Broadcast:
			h.mu.RLock()
			// Send to Recipient
			if conns, ok := h.Clients[msg.ToUserID]; ok {
				for conn := range conns {
					conn.WriteJSON(msg)
				}
			}
			// Send back to Sender (so their other tabs update, or for confirmation)
			if conns, ok := h.Clients[msg.FromUserID]; ok {
				for conn := range conns {
					conn.WriteJSON(msg)
				}
			}
			h.mu.RUnlock()
		}
	}
}

// SendToUser is a helper method to push a message into the hub
func (h *ChatHub) SendToUser(fromID, toID uint, content, timeStr string) {
	h.Broadcast <- &ChatMessage{
		FromUserID: fromID,
		ToUserID:   toID,
		Content:    content,
		Time:       timeStr,
	}
}
