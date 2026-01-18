package services

import (
	"FitClassMaster/internal/config"
	"FitClassMaster/internal/models"
	"FitClassMaster/internal/repositories"
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
