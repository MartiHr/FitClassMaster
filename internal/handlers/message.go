package handlers

import (
	"FitClassMaster/internal/auth"
	"FitClassMaster/internal/services"
	"FitClassMaster/internal/templates"
	"FitClassMaster/internal/websockets" // Make sure this import path is correct
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

type MessageHandler struct {
	Service     *services.MessageService
	UserService *services.UserService
	ChatHub     *websockets.ChatHub // <--- 1. We added the Hub here
}

// NewMessageHandler now requires the ChatHub as an argument
func NewMessageHandler(s *services.MessageService, us *services.UserService, hub *websockets.ChatHub) *MessageHandler {
	return &MessageHandler{
		Service:     s,
		UserService: us,
		ChatHub:     hub,
	}
}

// --- STANDARD PAGES (Inbox & Thread) ---

func (h *MessageHandler) InboxPage(w http.ResponseWriter, r *http.Request) {
	userID, _ := auth.GetUserIDFromSession(r)

	inbox, err := h.Service.GetInbox(userID)
	if err != nil {
		http.Error(w, "Failed to load inbox", http.StatusInternalServerError)
		return
	}

	data := map[string]any{
		"Title":      "Messages",
		"Inbox":      inbox,
		"ActiveConv": uint(0),
	}
	templates.SmartRender(w, r, "messages", "", data)
}

func (h *MessageHandler) ThreadPage(w http.ResponseWriter, r *http.Request) {
	userID, _ := auth.GetUserIDFromSession(r)
	convIDStr := chi.URLParam(r, "id")
	convID, _ := strconv.ParseUint(convIDStr, 10, 32)

	inbox, _ := h.Service.GetInbox(userID)
	messages, err := h.Service.GetThread(uint(convID))

	// Find Recipient ID for the UI
	var recipientID uint
	for _, chat := range inbox {
		if chat.ConversationID == uint(convID) {
			recipientID = chat.OtherUser.ID
			break
		}
	}

	if err != nil {
		http.Error(w, "Failed to load messages", http.StatusInternalServerError)
		return
	}

	data := map[string]any{
		"Title":         "Chat",
		"Inbox":         inbox,
		"ActiveConv":    uint(convID),
		"Messages":      messages,
		"CurrentUserID": userID,
		"RecipientID":   recipientID,
	}
	templates.SmartRender(w, r, "messages", "", data)
}

// StartChat handles the "Message Me" button click
func (h *MessageHandler) StartChat(w http.ResponseWriter, r *http.Request) {
	userID, _ := auth.GetUserIDFromSession(r)
	targetIDStr := chi.URLParam(r, "userID")
	targetID, _ := strconv.ParseUint(targetIDStr, 10, 32)

	conv, err := h.Service.StartConversation(userID, uint(targetID))
	if err != nil {
		http.Error(w, "Failed to start chat", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/messages/"+strconv.Itoa(int(conv.ID)), http.StatusSeeOther)
}

// --- SENDING LOGIC (Database + WebSocket) ---

func (h *MessageHandler) SendPost(w http.ResponseWriter, r *http.Request) {
	userID, _ := auth.GetUserIDFromSession(r)
	recipientID, _ := strconv.ParseUint(r.FormValue("recipient_id"), 10, 32)
	content := r.FormValue("content")

	if content == "" {
		http.Redirect(w, r, "/messages", http.StatusSeeOther)
		return
	}

	// Save to Database (So it's there next time you login)
	err := h.Service.SendMessage(userID, uint(recipientID), content)
	if err != nil {
		http.Error(w, "Failed to send", http.StatusInternalServerError)
		return
	}

	// Broadcast via WebSocket (So it appears instantly)
	// We send "time.Now()" just for display purposes
	h.ChatHub.SendToUser(userID, uint(recipientID), content, time.Now().Format("15:04"))

	// Response
	// If the request came from JavaScript (Fetch), we just return OK
	if r.Header.Get("HX-Request") != "" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Fallback for non-JS browsers
	http.Redirect(w, r, fmt.Sprintf("/messages"), http.StatusSeeOther)
}

// ServeWS is the endpoint: ws://localhost:8080/ws/chat
func (h *MessageHandler) ServeWS(w http.ResponseWriter, r *http.Request) {
	// Check Auth
	userID, ok := auth.GetUserIDFromSession(r)
	if !ok {
		// If not logged in, ignore the connection attempt
		return
	}

	// Upgrade the HTTP connection to a WebSocket connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	// Register the user in the ChatHub
	client := &websockets.ChatClient{UserID: userID, Conn: conn}
	h.ChatHub.Register <- client

	// Keep the connection open
	// This loop blocks until the user closes the tab
	defer func() {
		h.ChatHub.Unregister <- client
	}()

	for {
		// We read from the client just to keep the connection alive (Ping/Pong)
		// We ignore incoming messages here because we use the POST form for sending
		if _, _, err := conn.NextReader(); err != nil {
			break
		}
	}
}
