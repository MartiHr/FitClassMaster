package models

import "time"

type Class struct {
	ID          uint   `gorm:"primaryKey"`
	Name        string `gorm:"not null"`
	Description string
	TrainerID   uint
	Trainer     User      `gorm:"foreignKey:TrainerID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	StartTime   time.Time `gorm:"not null"`
	Duration    int       `gorm:"not null"` // Minutes
	MaxCapacity int       `gorm:"not null"`
	CreatedAt   time.Time
}
