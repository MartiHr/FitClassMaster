package repositories

import (
	"FitClassMaster/internal/config"
	"FitClassMaster/internal/models"
)

// EnrollmentRepo handles database operations for the Enrollment model.
type EnrollmentRepo struct{}

// NewEnrollmentRepo creates a new instance of EnrollmentRepo.
func NewEnrollmentRepo() *EnrollmentRepo {
	return &EnrollmentRepo{}
}

// Create inserts a new enrollment record.
func (r *EnrollmentRepo) Create(enrollment *models.Enrollment) error {
	return config.DB.Create(enrollment).Error
}

// Exists checks if an active enrollment already exists for a given user and class.
func (r *EnrollmentRepo) Exists(userID, classID uint) (bool, error) {
	var count int64
	err := config.DB.Model(&models.Enrollment{}).
		Where("user_id = ? AND class_id = ? AND status = 'active'", userID, classID).
		Count(&count).Error
	return count > 0, err
}

// CountActive returns the number of active enrollments for a specific class.
func (r *EnrollmentRepo) CountActive(classID uint) (int64, error) {
	var count int64
	err := config.DB.Model(&models.Enrollment{}).
		Where("class_id = ? AND status = 'active'", classID).
		Count(&count).Error
	return count, err
}

// GetUserActiveEnrollments retrieves all active enrollments for a specific user, including preloaded class details.
func (r *EnrollmentRepo) GetUserActiveEnrollments(userID uint) ([]models.Enrollment, error) {
	var enrollments []models.Enrollment
	// Preload "Class" so name and time data are available for the UI.
	err := config.DB.Preload("Class").
		Where("user_id = ? AND status = 'active'", userID).
		Find(&enrollments).Error
	return enrollments, err
}

// Delete removes an enrollment record (hard delete).
func (r *EnrollmentRepo) Delete(userID, classID uint) error {
	return config.DB.Where("user_id = ? AND class_id = ?", userID, classID).
		Delete(&models.Enrollment{}).Error
}
