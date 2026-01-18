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
	// 1. Get Class Details via ClassRepo to check MaxCapacity
	class, err := s.ClassRepo.GetByID(classID)
	if err != nil {
		return errors.New("class not found")
	}

	// 2. Check current capacity via EnrollmentRepo
	currentCount, _ := s.Repo.CountActive(classID)
	if int(currentCount) >= class.MaxCapacity {
		return errors.New("this class is full")
	}

	// 3. Check if already enrolled (using your existing Repo.Exists)
	exists, _ := s.Repo.Exists(userID, classID)
	if exists {
		return errors.New("you are already enrolled in this class")
	}

	// 4. Create enrollment
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
