package models

import "time"

type Class struct {
	ID          uint   `gorm:"primaryKey"`
	Name        string `gorm:"not null"`
	Description string

	StartTime time.Time `gorm:"not null"`
	Duration  int       `gorm:"not null"` // Minutes
	CreatedAt time.Time

	DifficultyLevel string `json:"difficulty_level"`

	MaxCapacity int `gorm:"not null"`

	TrainerID uint
	Trainer   User `gorm:"foreignKey:TrainerID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`

	Enrollments []Enrollment `gorm:"foreignKey:ClassID" json:"enrollments"`
}
