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

// GetActiveSessionsForTrainer returns in_progress sessions using plans created by this trainer
func (r *SessionRepo) GetActiveSessionsForTrainer(trainerID uint) ([]models.WorkoutSession, error) {
	var sessions []models.WorkoutSession
	err := config.DB.
		Preload("User").
		Preload("WorkoutPlan").
		Joins("JOIN workout_plans ON workout_plans.id = workout_sessions.workout_plan_id").
		Where("workout_sessions.status = ? AND workout_plans.trainer_id = ?", "in_progress", trainerID).
		Order("workout_sessions.start_time desc").
		Find(&sessions).Error
	return sessions, err
}

// GetActiveSessions fetches all sessions with status 'in_progress'
// Preloads User info so we know who is working out
// Used for admin
func (r *SessionRepo) GetActiveSessions() ([]models.WorkoutSession, error) {
	var sessions []models.WorkoutSession
	err := config.DB.
		Preload("User"). // To show "John Doe is working out"
		Preload("WorkoutPlan"). // To show "Leg Day"
		Where("status = ?", "in_progress").
		Order("start_time desc").
		Find(&sessions).Error
	return sessions, err
}

// GetCompletedByUserID fetches finished sessions for a specific member history
func (r *SessionRepo) GetCompletedByUserID(userID uint) ([]models.WorkoutSession, error) {
	var sessions []models.WorkoutSession
	err := config.DB.
		Preload("WorkoutPlan").
		Where("user_id = ? AND status = ?", userID, "completed").
		Order("start_time desc").
		Find(&sessions).Error
	return sessions, err
}

// GetRecentCompletedForTrainer fetches the last 5 completed sessions for a trainer's clients
func (r *SessionRepo) GetRecentCompletedForTrainer(trainerID uint) ([]models.WorkoutSession, error) {
	var sessions []models.WorkoutSession
	err := config.DB.
		Preload("User").
		Preload("WorkoutPlan").
		Joins("JOIN workout_plans ON workout_plans.id = workout_sessions.workout_plan_id").
		Where("workout_sessions.status = ? AND workout_plans.trainer_id = ?", "completed", trainerID).
		Order("workout_sessions.end_time desc").
		Limit(5).
		Find(&sessions).Error
	return sessions, err
}

func (r *SessionRepo) GetActiveByUserID(userID uint) (*models.WorkoutSession, error) {
	var session models.WorkoutSession

	// We look for status 'in_progress'
	err := config.DB.
		Preload("WorkoutPlan").
		Where("user_id = ? AND status = ?", userID, "in_progress").
		First(&session).Error

	return &session, err
}
