package services

import (
	"FitClassMaster/internal/models"
	"FitClassMaster/internal/repositories"
	"time"
)

type SessionService struct {
	Repo *repositories.SessionRepo
}

func NewSessionService(repo *repositories.SessionRepo) *SessionService {
	return &SessionService{Repo: repo}
}

// StartSession initializes a new workout session for a user based on a plan
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

func (s *SessionService) GetDetails(id uint) (*models.WorkoutSession, error) {
	return s.Repo.GetByID(id)
}

// LogSet validates (if needed) and saves the session log
func (s *SessionService) LogSet(logEntry *models.SessionLog) error {
	return s.Repo.SaveLog(logEntry)
}

// ListActiveSessions decides whether to fetch ALL or filtered based on role/ID
func (s *SessionService) ListActiveSessions(role models.Role, userID uint) ([]models.WorkoutSession, error) {
	if role == models.RoleAdmin {
		// Admins see everything
		return s.Repo.GetActiveSessions()
	}

	// Trainers only see sessions using THEIR plans
	return s.Repo.GetActiveSessionsForTrainer(userID)
}

// FinishSession marks the session as completed and sets the end time
func (s *SessionService) FinishSession(id uint) error {
	session, err := s.Repo.GetByID(id)
	if err != nil {
		return err
	}

	session.Status = "completed"
	session.EndTime = time.Now()

	return s.Repo.Update(session)
}

func (s *SessionService) GetHistory(userID uint) ([]models.WorkoutSession, error) {
	return s.Repo.GetCompletedByUserID(userID)
}

func (s *SessionService) GetTrainerHistory(trainerID uint) ([]models.WorkoutSession, error) {
	return s.Repo.GetRecentCompletedForTrainer(trainerID)
}

// GetActiveSession checks if the specific user has an unfinished workout
func (s *SessionService) GetActiveSession(userID uint) (*models.WorkoutSession, error) {
	return s.Repo.GetActiveByUserID(userID)
}
