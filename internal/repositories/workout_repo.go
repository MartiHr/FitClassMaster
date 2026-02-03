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
			return db.Order("workout_exercises.[order] asc") // Ensure correct order
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
// internal/repositories/workout_repo.go
func (r *WorkoutRepo) GetAll() ([]models.WorkoutPlan, error) {
	var plans []models.WorkoutPlan

	err := config.DB.
		Preload("Trainer").
		Preload("WorkoutExercises").
		Find(&plans).Error

	return plans, err
}

// ClearExercises removes all exercise rows for a specific plan
// Used when editing a plan (we delete old rows and re-insert new ones)
func (r *WorkoutRepo) ClearExercises(planID uint) error {
	return config.DB.Unscoped().Where("workout_plan_id = ?", planID).Delete(&models.WorkoutExercise{}).Error
}

// Update saves changes to the plan (Name, Notes, and new Exercises list)
func (r *WorkoutRepo) Update(plan *models.WorkoutPlan) error {
	return config.DB.Save(plan).Error
}

// Delete performs a Soft Delete (Update) so FK constraints don't break
func (r *WorkoutRepo) Delete(id uint) error {
	return config.DB.Delete(&models.WorkoutPlan{}, id).Error
}
