package models

import (
	"time"

	"gorm.io/gorm"
)

// Conversation represents a chat history between two users.
type Conversation struct {
	gorm.Model
	// User1ID and User2ID represent the participants.
	// Conventionally, User1ID < User2ID to ensure a unique conversation pair.
	User1ID uint `gorm:"not null;index"`
	User1   User `gorm:"foreignKey:User1ID"`

	User2ID uint `gorm:"not null;index"`
	User2   User `gorm:"foreignKey:User2ID"`

	LastMessageAt time.Time `gorm:"index"` // Used to sort conversations by the most recent activity.

	Messages []Message `gorm:"foreignKey:ConversationID"`
}

// Message represents a single text message sent within a conversation.
type Message struct {
	gorm.Model
	ConversationID uint         `gorm:"not null;index"`
	Conversation   Conversation `gorm:"foreignKey:ConversationID"`

	SenderID uint `gorm:"not null"`
	Sender   User `gorm:"foreignKey:SenderID"`

	Content string `gorm:"type:text;not null"`
	IsRead  bool   `gorm:"default:false"`
}
