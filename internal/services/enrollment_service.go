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

func (s *EnrollmentService) IsUserEnrolled(userID, classID uint) (bool, error) {
	return s.Repo.Exists(userID, classID)
}

func (s *EnrollmentService) EnrollUser(userID, classID uint) error {
	// Get Class Details via ClassRepo to check MaxCapacity
	class, err := s.ClassRepo.GetByID(classID)
	if err != nil {
		return errors.New("class not found")
	}

	// Check current capacity via EnrollmentRepo
	currentCount, _ := s.Repo.CountActive(classID)
	if int(currentCount) >= class.MaxCapacity {
		return errors.New("this class is full")
	}

	// Check if already enrolled
	exists, _ := s.IsUserEnrolled(userID, classID)

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

func (s *EnrollmentService) GetMySchedule(userID uint) ([]models.Class, error) {
	enrollments, err := s.Repo.GetUserActiveEnrollments(userID)
	if err != nil {
		return nil, err
	}

	var classes []models.Class
	for _, e := range enrollments {
		classes = append(classes, e.Class)
	}

	return classes, nil
}

func (s *EnrollmentService) CancelEnrollment(userID, classID uint) error {
	return s.Repo.Delete(userID, classID)
}
