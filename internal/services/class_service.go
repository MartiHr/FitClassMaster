package services

import (
	"FitClassMaster/internal/config"
	"FitClassMaster/internal/models"
	"FitClassMaster/internal/repositories"
	"time"
)

type ClassService struct {
	Repo *repositories.ClassRepo
}

func NewClassService(repo *repositories.ClassRepo) *ClassService {
	return &ClassService{Repo: repo}
}

func (s *ClassService) GetAvailableClasses() ([]models.Class, error) {
	return s.Repo.GetAll()
}

func (s *ClassService) GetClassesForUser(userID uint) ([]map[string]interface{}, error) {
	var classes []models.Class
	config.DB.Find(&classes) // Fetch all classes

	var userEnrollments []uint
	config.DB.Model(&models.Enrollment{}).
		Where("user_id = ? AND status = 'active'", userID).
		Pluck("class_id", &userEnrollments) // Get IDs of classes user is in

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

func (s *ClassService) GetFullDetails(classID uint) (*models.Class, error) {
	return s.Repo.GetWithDetails(classID)
}

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

func (s *ClassService) CancelClass(id uint) error {
	return s.Repo.Delete(id)
}

// UpdateClass modifies an existing class
func (s *ClassService) UpdateClass(id uint, name, description, difficulty string, start time.Time, durationMinutes int, capacity int) error {
	// Get the existing class
	class, err := s.Repo.GetByID(id)
	if err != nil {
		return err
	}

	// Update fields
	class.Name = name
	class.Description = description
	class.DifficultyLevel = difficulty
	class.StartTime = start
	class.Duration = durationMinutes
	class.MaxCapacity = capacity

	return s.Repo.Update(class)
}
