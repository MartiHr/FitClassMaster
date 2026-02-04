package services

import (
	"FitClassMaster/internal/models"
	"FitClassMaster/internal/repositories"
)

type MessageService struct {
	Repo *repositories.MessageRepo
}

func NewMessageService(repo *repositories.MessageRepo) *MessageService {
	return &MessageService{Repo: repo}
}

// SendMessage handles finding the chat ID and saving the text
func (s *MessageService) SendMessage(senderID, recipientID uint, content string) error {
	// Find or start conversation
	conv, err := s.Repo.FindOrCreateConversation(senderID, recipientID)
	if err != nil {
		return err
	}

	// Create Message object
	msg := &models.Message{
		ConversationID: conv.ID,
		SenderID:       senderID,
		Content:        content,
	}

	// Save to DB
	return s.Repo.CreateMessage(msg)
}

// ChatSummary is a helper struct for the UI (Inbox List)
type ChatSummary struct {
	ConversationID uint
	OtherUser      models.User
	LastMessage    string
	LastTime       string
}

// GetInbox prepares the data for the sidebar
func (s *MessageService) GetInbox(myUserID uint) ([]ChatSummary, error) {
	rawConvs, err := s.Repo.GetUserConversations(myUserID)
	if err != nil {
		return nil, err
	}

	var inbox []ChatSummary
	for _, c := range rawConvs {
		// Determine who the "Other" person is
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

// GetThread gets the full history for the active chat window
func (s *MessageService) GetThread(convID uint) ([]models.Message, error) {
	return s.Repo.GetHistory(convID)
}

// StartConversation gets the conversation ID without sending a message
func (s *MessageService) StartConversation(u1, u2 uint) (*models.Conversation, error) {
	return s.Repo.FindOrCreateConversation(u1, u2)
}
