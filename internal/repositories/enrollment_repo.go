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

func (r *EnrollmentRepo) CountActive(classID uint) (int64, error) {
	var count int64
	err := config.DB.Model(&models.Enrollment{}).
		Where("class_id = ? AND status = 'active'", classID).
		Count(&count).Error
	return count, err
}

func (r *EnrollmentRepo) GetUserActiveEnrollments(userID uint) ([]models.Enrollment, error) {
	var enrollments []models.Enrollment
	// Preload("Class") so we have the name/time data for the dashboard
	err := config.DB.Preload("Class").
		Where("user_id = ? AND status = 'active'", userID).
		Find(&enrollments).Error
	return enrollments, err
}

func (r *EnrollmentRepo) Delete(userID, classID uint) error {
	// Soft delete the enrollment record
	return config.DB.Where("user_id = ? AND class_id = ?", userID, classID).
		Delete(&models.Enrollment{}).Error
}
