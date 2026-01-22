package repositories

import (
	"FitClassMaster/internal/config"
	"FitClassMaster/internal/models"

	"gorm.io/gorm"
)

type SessionRepo struct{}

func NewSessionRepo() *SessionRepo {
	return &SessionRepo{}
}

// Create starts a new session
func (r *SessionRepo) Create(session *models.WorkoutSession) error {
	return config.DB.Create(session).Error
}

// GetByID fetches a session with all the plan details needed to display the UI
func (r *SessionRepo) GetByID(id uint) (*models.WorkoutSession, error) {
	var session models.WorkoutSession
	err := config.DB.
		Preload("WorkoutPlan").
		Preload("WorkoutPlan.WorkoutExercises", func(db *gorm.DB) *gorm.DB {
			return db.Order("workout_exercises.[order] asc")
		}).
		Preload("WorkoutPlan.WorkoutExercises.Exercise").
		Preload("Logs").
		First(&session, id).Error
	return &session, err
}

// Update saves changes (e.g., changing status to 'completed')
func (r *SessionRepo) Update(session *models.WorkoutSession) error {
	return config.DB.Save(session).Error
}

// SaveLog records a single set (reps/weight) to the database
func (r *SessionRepo) SaveLog(log *models.SessionLog) error {
	return config.DB.Create(log).Error
}
