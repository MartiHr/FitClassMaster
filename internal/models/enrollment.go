package models

import "time"

type Enrollment struct {
	ID        uint   `gorm:"primaryKey"`
	UserID    uint   `gorm:"not null"`
	ClassID   uint   `gorm:"not null"`
	Status    string `gorm:"size:20;default:'active'"` // active, cancelled
	CreatedAt time.Time

	// Relationships for Preloading
	User  User  `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;"`
	Class Class `gorm:"foreignKey:ClassID;constraint:OnDelete:CASCADE;"`
}
