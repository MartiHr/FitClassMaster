package repositories

import (
	"FitClassMaster/internal/config"
	"FitClassMaster/internal/models"
)

type EnrollmentRepo struct{}

func NewEnrollmentRepo() *EnrollmentRepo {
	return &EnrollmentRepo{}
}

func (r *EnrollmentRepo) Create(enrollment *models.Enrollment) error {
	return config.DB.Create(enrollment).Error
}

func (r *EnrollmentRepo) Exists(userID, classID uint) (bool, error) {
	var count int64
	err := config.DB.Model(&models.Enrollment{}).
		Where("user_id = ? AND class_id = ? AND status = 'active'", userID, classID).
		Count(&count).Error
	return count > 0, err
}
