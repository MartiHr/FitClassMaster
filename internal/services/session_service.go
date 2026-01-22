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
