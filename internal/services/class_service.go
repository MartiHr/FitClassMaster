package services

import (
	"FitClassMaster/internal/config"
	"FitClassMaster/internal/models"
	"FitClassMaster/internal/repositories"
	"time"
)

// ClassService manages fitness class schedules and user enrollment status.
type ClassService struct {
	Repo *repositories.ClassRepo
}

// NewClassService creates a new instance of ClassService.
func NewClassService(repo *repositories.ClassRepo) *ClassService {
	return &ClassService{Repo: repo}
}

// GetAvailableClasses retrieves all scheduled classes.
func (s *ClassService) GetAvailableClasses() ([]models.Class, error) {
	return s.Repo.GetAll()
}

// GetClassesForUser retrieves all classes and indicates if the specified user is enrolled in each.
func (s *ClassService) GetClassesForUser(userID uint) ([]map[string]interface{}, error) {
	var classes []models.Class
	if err := config.DB.Find(&classes).Error; err != nil {
		return nil, err
	}

	var userEnrollments []uint
	if err := config.DB.Model(&models.Enrollment{}).
		Where("user_id = ? AND status = 'active'", userID).
		Pluck("class_id", &userEnrollments).Error; err != nil {
		return nil, err
	}

	var results []map[string]interface{}
	for _, c := range classes {
		isEnrolled := false
		for _, id := range userEnrollments {
			if id == c.ID {
				isEnrolled = true
				break
			}
		}
		results = append(results, map[string]interface{}{
			"Class":      c,
			"IsEnrolled": isEnrolled,
		})
	}

	return results, nil
}

// GetFullDetails retrieves a class with preloaded trainer and participant information.
func (s *ClassService) GetFullDetails(classID uint) (*models.Class, error) {
	return s.Repo.GetWithDetails(classID)
}

// CreateClass initializes and saves a new class schedule.
func (s *ClassService) CreateClass(name, description, difficulty string, trainerID uint, start time.Time, durationMinutes int, capacity int) error {
	class := &models.Class{
		Name:            name,
		Description:     description,
		DifficultyLevel: difficulty,
		TrainerID:       trainerID,
		StartTime:       start,
		Duration:        durationMinutes,
		MaxCapacity:     capacity,
	}
	return s.Repo.Create(class)
}

// CancelClass removes a class from the schedule.
func (s *ClassService) CancelClass(id uint) error {
	return s.Repo.Delete(id)
}

// UpdateClass modifies an existing class schedule.
func (s *ClassService) UpdateClass(id uint, name, description, difficulty string, start time.Time, durationMinutes int, capacity int) error {
	class, err := s.Repo.GetByID(id)
	if err != nil {
		return err
	}

	class.Name = name
	class.Description = description
	class.DifficultyLevel = difficulty
	class.StartTime = start
	class.Duration = durationMinutes
	class.MaxCapacity = capacity

	return s.Repo.Update(class)
}
