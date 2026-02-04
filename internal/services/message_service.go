package services

import (
	"FitClassMaster/internal/models"
	"FitClassMaster/internal/repositories"
)

// MessageService manages real-time chat and message persistence logic.
type MessageService struct {
	Repo *repositories.MessageRepo
}

// NewMessageService creates a new instance of MessageService.
func NewMessageService(repo *repositories.MessageRepo) *MessageService {
	return &MessageService{Repo: repo}
}

// SendMessage handles finding or creating a conversation and persisting a new message.
func (s *MessageService) SendMessage(senderID, recipientID uint, content string) error {
	// Find or start a conversation between the two participants.
	conv, err := s.Repo.FindOrCreateConversation(senderID, recipientID)
	if err != nil {
		return err
	}

	// Initialize the message object.
	msg := &models.Message{
		ConversationID: conv.ID,
		SenderID:       senderID,
		Content:        content,
	}

	// Persist the message to the database.
	return s.Repo.CreateMessage(msg)
}

// ChatSummary is a DTO for displaying a conversation preview in the inbox.
type ChatSummary struct {
	ConversationID uint
	OtherUser      models.User
	LastMessage    string
	LastTime       string
}

// GetInbox prepares conversation summaries for a specific user's inbox view.
func (s *MessageService) GetInbox(myUserID uint) ([]ChatSummary, error) {
	rawConvs, err := s.Repo.GetUserConversations(myUserID)
	if err != nil {
		return nil, err
	}

	var inbox []ChatSummary
	for _, c := range rawConvs {
		// Identify the other participant in the conversation.
		other := c.User1
		if c.User1ID == myUserID {
			other = c.User2
		}

		preview := "No messages yet"
		if len(c.Messages) > 0 {
			preview = c.Messages[0].Content
		}

		inbox = append(inbox, ChatSummary{
			ConversationID: c.ID,
			OtherUser:      other,
			LastMessage:    preview,
			LastTime:       c.LastMessageAt.Format("Jan 02"),
		})
	}
	return inbox, nil
}

// GetThread retrieves the full chronological message history for a conversation.
func (s *MessageService) GetThread(convID uint) ([]models.Message, error) {
	return s.Repo.GetHistory(convID)
}

// StartConversation retrieves the conversation ID between two users without sending an initial message.
func (s *MessageService) StartConversation(u1, u2 uint) (*models.Conversation, error) {
	return s.Repo.FindOrCreateConversation(u1, u2)
}
