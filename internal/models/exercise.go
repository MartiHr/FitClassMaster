package models

import (
	"gorm.io/gorm"
)

// Exercise represents a single physical activity or movement.
type Exercise struct {
	gorm.Model
	Name        string `gorm:"size:100;not null" json:"name"`
	Description string `gorm:"type:text" json:"description"`

	// Categorization for filtering.
	MuscleGroup string `gorm:"size:50" json:"muscle_group"` // e.g., "Chest", "Legs"
	Equipment   string `gorm:"size:100" json:"equipment"`   // e.g., "Dumbbell", "None"

	// VideoURL is an optional link to a demo video for the exercise.
	VideoURL string `gorm:"size:255" json:"video_url"`
}
