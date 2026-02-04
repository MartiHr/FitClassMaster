package services

import (
	"FitClassMaster/internal/models"
	"FitClassMaster/internal/repositories"
	"errors"
)

// EnrollmentService manages the process of users signing up for and canceling classes.
type EnrollmentService struct {
	Repo      *repositories.EnrollmentRepo
	ClassRepo *repositories.ClassRepo
}

// NewEnrollmentService creates a new instance of EnrollmentService.
func NewEnrollmentService(repo *repositories.EnrollmentRepo, classRepo *repositories.ClassRepo) *EnrollmentService {
	return &EnrollmentService{Repo: repo, ClassRepo: classRepo}
}

// IsUserEnrolled checks if a specific user has an active enrollment in a class.
func (s *EnrollmentService) IsUserEnrolled(userID, classID uint) (bool, error) {
	return s.Repo.Exists(userID, classID)
}

// EnrollUser attempts to register a user for a class, enforcing capacity and duplicate registration rules.
func (s *EnrollmentService) EnrollUser(userID, classID uint) error {
	// Retrieve class details to check maximum capacity.
	class, err := s.ClassRepo.GetByID(classID)
	if err != nil {
		return errors.New("class not found")
	}

	// Count current active enrollments for the class.
	currentCount, _ := s.Repo.CountActive(classID)
	if int(currentCount) >= class.MaxCapacity {
		return errors.New("this class is full")
	}

	// Prevent duplicate enrollments.
	exists, _ := s.IsUserEnrolled(userID, classID)
	if exists {
		return errors.New("you are already enrolled in this class")
	}

	// Create a new enrollment record.
	enrollment := &models.Enrollment{
		UserID:  userID,
		ClassID: classID,
		Status:  "active",
	}
	return s.Repo.Create(enrollment)
}

// GetMySchedule retrieves all classes that the specified user is currently enrolled in.
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

// CancelEnrollment removes a user's enrollment from a specific class.
func (s *EnrollmentService) CancelEnrollment(userID, classID uint) error {
	return s.Repo.Delete(userID, classID)
}
