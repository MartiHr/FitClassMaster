package repositories

import (
	"FitClassMaster/internal/config"
	"FitClassMaster/internal/models"

	"gorm.io/gorm"
)

type WorkoutRepo struct{}

func NewWorkoutRepo() *WorkoutRepo {
	return &WorkoutRepo{}
}

// CreateWithExercises saves the plan AND all its exercise links in one transaction
func (r *WorkoutRepo) Create(plan *models.WorkoutPlan) error {
	return config.DB.Create(plan).Error
}

// GetByID fetches a plan and preloads the exercises + the specific sets/reps info
func (r *WorkoutRepo) GetByID(id uint) (*models.WorkoutPlan, error) {
	var plan models.WorkoutPlan
	err := config.DB.
		Preload("Trainer").
		// Preload the link table (WorkoutExercises)
		Preload("WorkoutExercises", func(db *gorm.DB) *gorm.DB {
			return db.Order("workout_exercises.order asc") // Ensure correct order
		}).
		// Preload the actual Exercise details inside the link table
		Preload("WorkoutExercises.Exercise").
		First(&plan, id).Error
	return &plan, err
}

// GetAllForTrainer fetches plans created by a specific trainer
func (r *WorkoutRepo) GetAllForTrainer(trainerID uint) ([]models.WorkoutPlan, error) {
	var plans []models.WorkoutPlan
	err := config.DB.Where("trainer_id = ?", trainerID).Find(&plans).Error
	return plans, err
}

// GetAll fetches all plans (for Members to browse if allowed, or Admins)
func (r *WorkoutRepo) GetAll() ([]models.WorkoutPlan, error) {
	var plans []models.WorkoutPlan
	err := config.DB.Preload("Trainer").Find(&plans).Error
	return plans, err
}
