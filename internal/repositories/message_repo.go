package repositories

import (
	"FitClassMaster/internal/config"
	"FitClassMaster/internal/models"
	"time"

	"gorm.io/gorm"
)

// MessageRepo handles database operations for Conversations and Messages.
type MessageRepo struct{}

// NewMessageRepo creates a new instance of MessageRepo.
func NewMessageRepo() *MessageRepo {
	return &MessageRepo{}
}

// FindOrCreateConversation ensures a chat conversation exists between two users.
func (r *MessageRepo) FindOrCreateConversation(u1ID, u2ID uint) (*models.Conversation, error) {
	var conv models.Conversation

	// Enforce order to prevent duplicate conversations (Smallest ID always first).
	first, second := u1ID, u2ID
	if u1ID > u2ID {
		first, second = u2ID, u1ID
	}

	// Try to find an existing conversation.
	err := config.DB.
		Where("user1_id = ? AND user2_id = ?", first, second).
		First(&conv).Error

	if err == nil {
		return &conv, nil
	}

	// If not found, create a new conversation record.
	newConv := models.Conversation{
		User1ID:       first,
		User2ID:       second,
		LastMessageAt: time.Now(),
	}
	err = config.DB.Create(&newConv).Error
	return &newConv, err
}

// GetUserConversations fetches all conversations involving a specific user, including preloaded participants and the latest message preview.
func (r *MessageRepo) GetUserConversations(userID uint) ([]models.Conversation, error) {
	var convs []models.Conversation
	err := config.DB.
		Preload("User1").
		Preload("User2").
		// Preview the last message for the inbox list.
		Preload("Messages", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at desc").Limit(1)
		}).
		Where("user1_id = ? OR user2_id = ?", userID, userID).
		Order("last_message_at desc").
		Find(&convs).Error
	return convs, err
}

// GetHistory fetches the full message history for a specific conversation, ordered by creation time.
func (r *MessageRepo) GetHistory(convID uint) ([]models.Message, error) {
	var msgs []models.Message
	err := config.DB.
		Preload("Sender").
		Where("conversation_id = ?", convID).
		Order("created_at asc"). // Chronological order.
		Find(&msgs).Error
	return msgs, err
}

// CreateMessage saves a new message and updates the last activity timestamp of the associated conversation within a transaction.
func (r *MessageRepo) CreateMessage(msg *models.Message) error {
	return config.DB.Transaction(func(tx *gorm.DB) error {
		// Save the message record.
		if err := tx.Create(msg).Error; err != nil {
			return err
		}

		// Update the conversation's last activity timestamp to move it to the top of the inbox.
		return tx.Model(&models.Conversation{}).
			Where("id = ?", msg.ConversationID).
			Update("last_message_at", msg.CreatedAt).Error
	})
}
