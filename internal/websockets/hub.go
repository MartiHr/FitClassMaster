package websockets

import (
	"FitClassMaster/internal/models"
	"FitClassMaster/internal/services"
	"log"
	"sync"
)

// WSMessage is the data moving through the hub
type WSMessage struct {
	Type       string  `json:"type"`
	SessionID  uint    `json:"session_id"`
	ExerciseID uint    `json:"exercise_id"`
	SetNumber  int     `json:"set_number"`
	Reps       int     `json:"reps"`
	Weight     float64 `json:"weight"`
}

type Hub struct {
	// Service needed to save data to DB
	SessionService *services.SessionService

	// Registered clients, grouped by Session ID (The "Room")
	// map[SessionID] -> map[Client] -> true
	rooms map[uint]map[*Client]bool

	// Lock to prevent race conditions when modifying the map
	mu sync.RWMutex

	// Inbound messages from the clients
	Broadcast chan WSMessage

	// Register requests from the clients
	Register chan *Client

	// Unregister requests from clients
	Unregister chan *Client
}

func NewHub(service *services.SessionService) *Hub {
	return &Hub{
		SessionService: service,
		Broadcast:      make(chan WSMessage),
		Register:       make(chan *Client),
		Unregister:     make(chan *Client),
		rooms:          make(map[uint]map[*Client]bool),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.mu.Lock()
			if h.rooms[client.SessionID] == nil {
				h.rooms[client.SessionID] = make(map[*Client]bool)
			}
			h.rooms[client.SessionID][client] = true
			h.mu.Unlock()
			log.Printf("🔌 Client joined Session Room %d", client.SessionID)

		case client := <-h.Unregister:
			h.mu.Lock()
			if _, ok := h.rooms[client.SessionID]; ok {
				delete(h.rooms[client.SessionID], client)
				close(client.Send)
				// Clean up empty rooms
				if len(h.rooms[client.SessionID]) == 0 {
					delete(h.rooms, client.SessionID)
				}
			}
			h.mu.Unlock()
			log.Printf("🔌 Client left Session Room %d", client.SessionID)

		case msg := <-h.Broadcast:
			// 1. Save to DB (Async)
			go func(m WSMessage) {
				logEntry := &models.SessionLog{
					SessionID:  m.SessionID,
					ExerciseID: m.ExerciseID,
					SetNumber:  m.SetNumber,
					Reps:       m.Reps,
					Weight:     m.Weight,
				}
				if err := h.SessionService.LogSet(logEntry); err != nil {
					log.Println("❌ Error saving log:", err)
				} else {
					log.Printf("💾 Saved Log for Session %d", m.SessionID)
				}
			}(msg)

			// 2. Broadcast to everyone in that specific room (Trainer monitoring)
			h.mu.RLock()
			if clients, ok := h.rooms[msg.SessionID]; ok {
				for client := range clients {
					select {
					case client.Send <- msg:
					default:
						close(client.Send)
						delete(clients, client)
					}
				}
			}
			h.mu.RUnlock()
		}
	}
}
