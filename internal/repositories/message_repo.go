package repositories

import (
	"FitClassMaster/internal/config"
	"FitClassMaster/internal/models"
	"time"

	"gorm.io/gorm"
)

type MessageRepo struct{}

func NewMessageRepo() *MessageRepo {
	return &MessageRepo{}
}

// FindOrCreateConversation ensures a chat exists between two users
func (r *MessageRepo) FindOrCreateConversation(u1ID, u2ID uint) (*models.Conversation, error) {
	var conv models.Conversation

	// Enforce order to prevent duplicates (Smallest ID always first)
	first, second := u1ID, u2ID
	if u1ID > u2ID {
		first, second = u2ID, u1ID
	}

	// Try to find existing conversation
	err := config.DB.
		Where("user1_id = ? AND user2_id = ?", first, second).
		First(&conv).Error

	if err == nil {
		return &conv, nil
	}

	// If not found, create new
	newConv := models.Conversation{
		User1ID:       first,
		User2ID:       second,
		LastMessageAt: time.Now(),
	}
	err = config.DB.Create(&newConv).Error
	return &newConv, err
}

// GetUserConversations fetches all chats for a user (Inbox view)
func (r *MessageRepo) GetUserConversations(userID uint) ([]models.Conversation, error) {
	var convs []models.Conversation
	err := config.DB.
		Preload("User1").
		Preload("User2").
		// Preview the last message for the inbox list
		Preload("Messages", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at desc").Limit(1)
		}).
		Where("user1_id = ? OR user2_id = ?", userID, userID).
		Order("last_message_at desc").
		Find(&convs).Error
	return convs, err
}

// GetHistory fetches the full message log for a specific chat
func (r *MessageRepo) GetHistory(convID uint) ([]models.Message, error) {
	var msgs []models.Message
	err := config.DB.
		Preload("Sender").
		Where("conversation_id = ?", convID).
		Order("created_at asc"). // Oldest at top
		Find(&msgs).Error
	return msgs, err
}

// CreateMessage saves the message AND updates the Conversation's timestamp
func (r *MessageRepo) CreateMessage(msg *models.Message) error {
	return config.DB.Transaction(func(tx *gorm.DB) error {
		// Save Message
		if err := tx.Create(msg).Error; err != nil {
			return err
		}

		// Bump the Conversation timestamp (so it moves to top of inbox)
		return tx.Model(&models.Conversation{}).
			Where("id = ?", msg.ConversationID).
			Update("last_message_at", msg.CreatedAt).Error
	})
}
