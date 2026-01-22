package repositories

import (
	"FitClassMaster/internal/config"
	"FitClassMaster/internal/models"
)

type ClassRepo struct{}

func NewClassRepo() *ClassRepo {
	return &ClassRepo{}
}

func (r *ClassRepo) GetAll() ([]models.Class, error) {
	var classes []models.Class

	err := config.DB.Order("start_time asc").Find(&classes).Error
	return classes, err
}

func (r *ClassRepo) GetByID(id uint) (*models.Class, error) {
	var class models.Class
	err := config.DB.First(&class, id).Error
	return &class, err
}

func (r *ClassRepo) GetWithDetails(id uint) (*models.Class, error) {
	var class models.Class
	// We use Preload to fetch related data in one go
	err := config.DB.Preload("Trainer").
		Preload("Enrollments.User").
		First(&class, id).Error

	return &class, err
}
