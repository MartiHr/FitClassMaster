package repositories

import (
	"FitClassMaster/internal/config"
	"FitClassMaster/internal/models"
)

// ExerciseRepo handles database operations for the Exercise model.
type ExerciseRepo struct{}

// NewExerciseRepo creates a new instance of ExerciseRepo.
func NewExerciseRepo() *ExerciseRepo {
	return &ExerciseRepo{}
}

// GetAll retrieves all exercises from the database.
func (r *ExerciseRepo) GetAll() ([]models.Exercise, error) {
	var exercises []models.Exercise
	err := config.DB.Find(&exercises).Error
	return exercises, err
}

// Create inserts a new exercise record into the database.
func (r *ExerciseRepo) Create(exercise *models.Exercise) error {
	return config.DB.Create(exercise).Error
}

// GetByID retrieves a single exercise by its primary key ID.
func (r *ExerciseRepo) GetByID(id uint) (*models.Exercise, error) {
	var exercise models.Exercise
	err := config.DB.First(&exercise, id).Error
	return &exercise, err
}

// Update saves changes to an existing exercise record.
func (r *ExerciseRepo) Update(exercise *models.Exercise) error {
	return config.DB.Save(exercise).Error
}

// Delete removes an exercise record from the database.
func (r *ExerciseRepo) Delete(id uint) error {
	return config.DB.Delete(&models.Exercise{}, id).Error
}
