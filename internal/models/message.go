package models

import (
	"time"

	"gorm.io/gorm"
)

// Conversation represents a chat history between two users
type Conversation struct {
	gorm.Model
	// To prevent duplicate conversations, we can enforce that User1ID < User2ID
	User1ID uint `gorm:"not null;index"`
	User1   User `gorm:"foreignKey:User1ID"`

	User2ID uint `gorm:"not null;index"`
	User2   User `gorm:"foreignKey:User2ID"`

	LastMessageAt time.Time `gorm:"index"` // To sort by "Recent"

	Messages []Message `gorm:"foreignKey:ConversationID"`
}

// Message is a single text sent from one user to another
type Message struct {
	gorm.Model
	ConversationID uint         `gorm:"not null;index"`
	Conversation   Conversation `gorm:"foreignKey:ConversationID"`

	SenderID uint `gorm:"not null"`
	Sender   User `gorm:"foreignKey:SenderID"`

	Content string `gorm:"type:text;not null"`
	IsRead  bool   `gorm:"default:false"`
}
