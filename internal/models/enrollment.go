package models

import "time"

// Enrollment links a User to a Class.
type Enrollment struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    uint      `gorm:"not null"`
	ClassID   uint      `gorm:"not null"`
	Status    string    `gorm:"size:20;default:'active'"` // Possible values: "active", "cancelled"
	CreatedAt time.Time `gorm:"autoCreateTime"`

	// Relationships for Preloading
	User  User  `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;"`
	Class Class `gorm:"foreignKey:ClassID;constraint:OnDelete:CASCADE;"`
}
