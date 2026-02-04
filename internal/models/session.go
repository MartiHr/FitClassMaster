package models

import (
	"time"

	"gorm.io/gorm"
)

// WorkoutSession represents an instance of a user performing a WorkoutPlan.
type WorkoutSession struct {
	gorm.Model
	UserID uint `json:"user_id"`
	User   User `gorm:"foreignKey:UserID"`

	WorkoutPlanID uint        `json:"workout_plan_id"`
	WorkoutPlan   WorkoutPlan `gorm:"foreignKey:WorkoutPlanID"`

	StartTime time.Time
	EndTime   time.Time
	Status    string `gorm:"default:'in_progress'"` // Possible values: 'in_progress', 'completed'

	// Logs contains the detailed records for each set performed during the session.
	Logs []SessionLog `gorm:"foreignKey:SessionID"`
}

// SessionLog records the actual performance of a single set for a specific exercise.
type SessionLog struct {
	gorm.Model
	SessionID uint `json:"session_id"`

	ExerciseID uint     `json:"exercise_id"`
	Exercise   Exercise `gorm:"foreignKey:ExerciseID"`

	SetNumber int     `json:"set_number"`
	Reps      int     `json:"reps"`
	Weight    float64 `json:"weight"` // Actual weight used by the user
	Notes     string  `json:"notes"`
}

// GetLog retrieves a specific SessionLog for a given exercise and set number from the session.
func (s *WorkoutSession) GetLog(exerciseID uint, setNum int) *SessionLog {
	for _, log := range s.Logs {
		if log.ExerciseID == exerciseID && log.SetNumber == setNum {
			return &log
		}
	}
	return nil
}
