package models

import (
	"time"

	"gorm.io/gorm"
)

type WorkoutSession struct {
	gorm.Model
	UserID uint `json:"user_id"`
	User   User `gorm:"foreignKey:UserID"`

	WorkoutPlanID uint        `json:"workout_plan_id"`
	WorkoutPlan   WorkoutPlan `gorm:"foreignKey:WorkoutPlanID"`

	StartTime time.Time
	EndTime   time.Time
	Status    string `gorm:"default:'in_progress'"` // 'in_progress' or 'completed'

	// One session has many logged sets
	Logs []SessionLog `gorm:"foreignKey:SessionID"`
}

type SessionLog struct {
	gorm.Model
	SessionID uint `json:"session_id"`

	ExerciseID uint     `json:"exercise_id"`
	Exercise   Exercise `gorm:"foreignKey:ExerciseID"`

	SetNumber int     `json:"set_number"`
	Reps      int     `json:"reps"`
	Weight    float64 `json:"weight"` // Users log the actual weight used
	Notes     string  `json:"notes"`
}

// GetLog finds a specific log for an exercise and set number
func (s *WorkoutSession) GetLog(exerciseID uint, setNum int) *SessionLog {
	for _, log := range s.Logs {
		if log.ExerciseID == exerciseID && log.SetNumber == setNum {
			return &log
		}
	}
	return nil
}
