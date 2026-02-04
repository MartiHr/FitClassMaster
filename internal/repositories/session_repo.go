package repositories

import (
	"FitClassMaster/internal/config"
	"FitClassMaster/internal/models"

	"gorm.io/gorm"
)

// SessionRepo handles database operations for WorkoutSessions and SessionLogs.
type SessionRepo struct{}

// NewSessionRepo creates a new instance of SessionRepo.
func NewSessionRepo() *SessionRepo {
	return &SessionRepo{}
}

// Create inserts a new workout session into the database.
func (r *SessionRepo) Create(session *models.WorkoutSession) error {
	return config.DB.Create(session).Error
}

// GetByID retrieves a workout session by its ID, preloading the plan, exercises, and logs.
func (r *SessionRepo) GetByID(id uint) (*models.WorkoutSession, error) {
	var session models.WorkoutSession
	err := config.DB.
		Preload("WorkoutPlan").
		Preload("WorkoutPlan.WorkoutExercises", func(db *gorm.DB) *gorm.DB {
			// Order exercises as defined in the plan.
			return db.Order("workout_exercises.[order] asc")
		}).
		Preload("WorkoutPlan.WorkoutExercises.Exercise").
		Preload("Logs").
		First(&session, id).Error
	return &session, err
}

// Update saves changes to an existing workout session (e.g., status update).
func (r *SessionRepo) Update(session *models.WorkoutSession) error {
	return config.DB.Save(session).Error
}

// SaveLog records a single set performance log.
func (r *SessionRepo) SaveLog(log *models.SessionLog) error {
	return config.DB.Create(log).Error
}

// GetActiveSessionsForTrainer retrieves all in-progress sessions for users following plans created by a specific trainer.
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

// GetActiveSessions retrieves all sessions currently in progress across the system.
func (r *SessionRepo) GetActiveSessions() ([]models.WorkoutSession, error) {
	var sessions []models.WorkoutSession
	err := config.DB.
		Preload("User").        // Preload user to show who is working out.
		Preload("WorkoutPlan"). // Preload plan to show what they are doing.
		Where("status = ?", "in_progress").
		Order("start_time desc").
		Find(&sessions).Error
	return sessions, err
}

// GetCompletedByUserID retrieves the finished workout history for a specific user.
func (r *SessionRepo) GetCompletedByUserID(userID uint) ([]models.WorkoutSession, error) {
	var sessions []models.WorkoutSession
	err := config.DB.
		Preload("WorkoutPlan").
		Where("user_id = ? AND status = ?", userID, "completed").
		Order("start_time desc").
		Find(&sessions).Error
	return sessions, err
}

// GetRecentCompletedForTrainer retrieves the 5 most recently completed sessions for users following a specific trainer's plans.
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

// GetActiveByUserID retrieves the current in-progress session for a specific user, if one exists.
func (r *SessionRepo) GetActiveByUserID(userID uint) (*models.WorkoutSession, error) {
	var session models.WorkoutSession
	err := config.DB.
		Preload("WorkoutPlan").
		Where("user_id = ? AND status = ?", userID, "in_progress").
		First(&session).Error

	return &session, err
}
