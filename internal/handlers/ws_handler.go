package handlers

import (
	"FitClassMaster/internal/websockets"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type WSHandler struct {
	Hub *websockets.Hub // <-- Now holds the Hub, not the Service directly
}

func NewWSHandler(hub *websockets.Hub) *WSHandler {
	return &WSHandler{Hub: hub}
}

func (h *WSHandler) HandleSessionConnection(w http.ResponseWriter, r *http.Request) {
	sessionIDStr := chi.URLParam(r, "id")
	sessionID, _ := strconv.ParseUint(sessionIDStr, 10, 32)

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}

	// Create the Client and register with Hub
	client := &websockets.Client{
		Hub:       h.Hub,
		Conn:      conn,
		SessionID: uint(sessionID),
		Send:      make(chan websockets.WSMessage, 256),
	}

	client.Hub.Register <- client

	// Start the background processes for this connection
	go client.WritePump()
	go client.ReadPump()
}
