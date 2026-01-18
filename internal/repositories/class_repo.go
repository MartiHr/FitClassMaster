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
