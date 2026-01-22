package repositories

import (
	"FitClassMaster/internal/config"
	"FitClassMaster/internal/models"
)

type ExerciseRepo struct{}

func NewExerciseRepo() *ExerciseRepo {
	return &ExerciseRepo{}
}

// GetAll fetches all exercises
func (r *ExerciseRepo) GetAll() ([]models.Exercise, error) {
	var exercises []models.Exercise
	err := config.DB.Find(&exercises).Error
	return exercises, err
}

// Create adds a new exercise
func (r *ExerciseRepo) Create(exercise *models.Exercise) error {
	return config.DB.Create(exercise).Error
}

// GetByID fetches a specific exercise
func (r *ExerciseRepo) GetByID(id uint) (*models.Exercise, error) {
	var exercise models.Exercise
	err := config.DB.First(&exercise, id).Error
	return &exercise, err
}

// Update modifies an existing exercise
func (r *ExerciseRepo) Update(exercise *models.Exercise) error {
	return config.DB.Save(exercise).Error
}

// Delete removes an exercise
func (r *ExerciseRepo) Delete(id uint) error {
	return config.DB.Delete(&models.Exercise{}, id).Error
}
