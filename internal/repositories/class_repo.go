package repositories

import (
	"FitClassMaster/internal/config"
	"FitClassMaster/internal/models"
)

// ClassRepo handles database operations for the Class model.
type ClassRepo struct{}

// NewClassRepo creates a new instance of ClassRepo.
func NewClassRepo() *ClassRepo {
	return &ClassRepo{}
}

// GetAll retrieves all classes, ordered by start time.
func (r *ClassRepo) GetAll() ([]models.Class, error) {
	var classes []models.Class
	err := config.DB.Order("start_time asc").Find(&classes).Error
	return classes, err
}

// GetByID retrieves a single class by its primary key ID.
func (r *ClassRepo) GetByID(id uint) (*models.Class, error) {
	var class models.Class
	err := config.DB.First(&class, id).Error
	return &class, err
}

// GetWithDetails retrieves a single class by ID including preloaded Trainer and Enrollment information.
func (r *ClassRepo) GetWithDetails(id uint) (*models.Class, error) {
	var class models.Class
	// We use Preload to fetch related data in one go.
	err := config.DB.Preload("Trainer").
		Preload("Enrollments.User").
		First(&class, id).Error

	return &class, err
}

// Create inserts a new class into the database.
func (r *ClassRepo) Create(class *models.Class) error {
	return config.DB.Create(class).Error
}

// Delete removes a class from the database.
func (r *ClassRepo) Delete(id uint) error {
	return config.DB.Delete(&models.Class{}, id).Error
}

// Update saves changes to an existing class record.
func (r *ClassRepo) Update(class *models.Class) error {
	return config.DB.Save(class).Error
}
