package services

import (
	"FitClassMaster/internal/models"
	"FitClassMaster/internal/repositories"
	"time"
)

// SessionService manages the lifecycle of workout sessions and performance logging.
type SessionService struct {
	Repo *repositories.SessionRepo
}

// NewSessionService creates a new instance of SessionService.
func NewSessionService(repo *repositories.SessionRepo) *SessionService {
	return &SessionService{Repo: repo}
}

// StartSession initializes a new workout session for a user based on a specific workout plan.
func (s *SessionService) StartSession(userID, planID uint) (*models.WorkoutSession, error) {
	session := &models.WorkoutSession{
		UserID:        userID,
		WorkoutPlanID: planID,
		StartTime:     time.Now(),
		Status:        "in_progress",
	}

	err := s.Repo.Create(session)
	return session, err
}

// GetDetails retrieves a workout session by ID with all related information.
func (s *SessionService) GetDetails(id uint) (*models.WorkoutSession, error) {
	return s.Repo.GetByID(id)
}

// LogSet saves a performance log entry for a single set.
func (s *SessionService) LogSet(logEntry *models.SessionLog) error {
	return s.Repo.SaveLog(logEntry)
}

// ListActiveSessions retrieves in-progress sessions filtered by the requester's role.
// Admins see all active sessions, while Trainers see sessions using their plans.
func (s *SessionService) ListActiveSessions(role models.Role, userID uint) ([]models.WorkoutSession, error) {
	if role == models.RoleAdmin {
		return s.Repo.GetActiveSessions()
	}

	return s.Repo.GetActiveSessionsForTrainer(userID)
}

// FinishSession marks a workout session as completed and records the end time.
func (s *SessionService) FinishSession(id uint) error {
	session, err := s.Repo.GetByID(id)
	if err != nil {
		return err
	}

	session.Status = "completed"
	session.EndTime = time.Now()

	return s.Repo.Update(session)
}

// GetHistory retrieves the completed workout session history for a specific user.
func (s *SessionService) GetHistory(userID uint) ([]models.WorkoutSession, error) {
	return s.Repo.GetCompletedByUserID(userID)
}

// GetTrainerHistory retrieves recent completed sessions for users following a specific trainer's plans.
func (s *SessionService) GetTrainerHistory(trainerID uint) ([]models.WorkoutSession, error) {
	return s.Repo.GetRecentCompletedForTrainer(trainerID)
}

// GetActiveSession retrieves the current active session for a specific user, if any.
func (s *SessionService) GetActiveSession(userID uint) (*models.WorkoutSession, error) {
	return s.Repo.GetActiveByUserID(userID)
}
