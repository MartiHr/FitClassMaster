package repositories

import (
	"FitClassMaster/internal/config"
	"FitClassMaster/internal/models"

	"gorm.io/gorm"
)

// WorkoutRepo handles database operations for WorkoutPlans and WorkoutExercises.
type WorkoutRepo struct{}

// NewWorkoutRepo creates a new instance of WorkoutRepo.
func NewWorkoutRepo() *WorkoutRepo {
	return &WorkoutRepo{}
}

// Create inserts a new workout plan and its associated exercises in one transaction.
func (r *WorkoutRepo) Create(plan *models.WorkoutPlan) error {
	return config.DB.Create(plan).Error
}

// GetByID retrieves a workout plan by its ID, preloading trainer and exercise details.
func (r *WorkoutRepo) GetByID(id uint) (*models.WorkoutPlan, error) {
	var plan models.WorkoutPlan
	err := config.DB.
		Preload("Trainer").
		// Preload the association table with ordering.
		Preload("WorkoutExercises", func(db *gorm.DB) *gorm.DB {
			return db.Order("workout_exercises.[order] asc")
		}).
		// Preload actual Exercise data.
		Preload("WorkoutExercises.Exercise").
		First(&plan, id).Error
	return &plan, err
}

// GetAllForTrainer retrieves all workout plans created by a specific trainer.
func (r *WorkoutRepo) GetAllForTrainer(trainerID uint) ([]models.WorkoutPlan, error) {
	var plans []models.WorkoutPlan
	err := config.DB.Where("trainer_id = ?", trainerID).Find(&plans).Error
	return plans, err
}

// GetAll retrieves all workout plans available in the system.
func (r *WorkoutRepo) GetAll() ([]models.WorkoutPlan, error) {
	var plans []models.WorkoutPlan
	err := config.DB.
		Preload("Trainer").
		Preload("WorkoutExercises").
		Find(&plans).Error

	return plans, err
}

// ClearExercises removes all associated exercises for a specific plan.
// Typically used during a plan update to refresh the exercise list.
func (r *WorkoutRepo) ClearExercises(planID uint) error {
	return config.DB.Unscoped().Where("workout_plan_id = ?", planID).Delete(&models.WorkoutExercise{}).Error
}

// Update saves changes to an existing workout plan record.
func (r *WorkoutRepo) Update(plan *models.WorkoutPlan) error {
	return config.DB.Save(plan).Error
}

// Delete performs a soft delete of a workout plan.
func (r *WorkoutRepo) Delete(id uint) error {
	return config.DB.Delete(&models.WorkoutPlan{}, id).Error
}
