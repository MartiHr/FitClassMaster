package models

import (
	"gorm.io/gorm"
)

type Exercise struct {
	gorm.Model
	Name        string `gorm:"size:100;not null" json:"name"`
	Description string `gorm:"type:text" json:"description"`

	// Categorization for filtering
	MuscleGroup string `gorm:"size:50" json:"muscle_group"`
	Equipment   string `gorm:"size:100" json:"equipment"`

	// Optional: Link to a demo video for the user to watch
	VideoURL string `gorm:"size:255" json:"video_url"`
}
