package models

import "time"

// Class represents a scheduled fitness class.
type Class struct {
	ID          uint   `gorm:"primaryKey"`
	Name        string `gorm:"not null"`
	Description string

	StartTime time.Time `gorm:"not null"`
	Duration  int       `gorm:"not null"` // Duration in minutes
	CreatedAt time.Time `gorm:"autoCreateTime"`

	DifficultyLevel string `json:"difficulty_level"` // e.g., "Beginner", "Intermediate", "Advanced"

	MaxCapacity int `gorm:"not null"`

	TrainerID uint
	Trainer   User `gorm:"foreignKey:TrainerID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`

	Enrollments []Enrollment `gorm:"foreignKey:ClassID" json:"enrollments"`
}
