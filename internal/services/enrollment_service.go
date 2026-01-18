package services

import (
	"FitClassMaster/internal/models"
	"FitClassMaster/internal/repositories"
	"errors"
)

type EnrollmentService struct {
	Repo      *repositories.EnrollmentRepo
	ClassRepo *repositories.ClassRepo
}

func NewEnrollmentService(repo *repositories.EnrollmentRepo, classRepo *repositories.ClassRepo) *EnrollmentService {
	return &EnrollmentService{Repo: repo, ClassRepo: classRepo}
}

func (s *EnrollmentService) EnrollUser(userID, classID uint) error {
	// Check if already enrolled
	exists, _ := s.Repo.Exists(userID, classID)
	if exists {
		return errors.New("you are already enrolled in this class")
	}

	// Create enrollment
	enrollment := &models.Enrollment{
		UserID:  userID,
		ClassID: classID,
		Status:  "active",
	}
	return s.Repo.Create(enrollment)
}
