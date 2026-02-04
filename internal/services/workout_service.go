package services

import (
	"FitClassMaster/internal/models"
	"FitClassMaster/internal/repositories"
	"strconv"
)

// WorkoutService manages the creation and updating of workout plans.
type WorkoutService struct {
	Repo *repositories.WorkoutRepo
}

// NewWorkoutService creates a new instance of WorkoutService.
func NewWorkoutService(repo *repositories.WorkoutRepo) *WorkoutService {
	return &WorkoutService{Repo: repo}
}

// ExerciseInput represents a single exercise row within a workout plan creation form.
type ExerciseInput struct {
	ExerciseID uint
	Sets       int
	Reps       int
	Notes      string
}

// CreatePlan initializes and persists a new workout plan with its associated exercises.
func (s *WorkoutService) CreatePlan(name, description string, trainerID uint, exercises []ExerciseInput) error {
	// Initialize the workout plan object.
	plan := models.WorkoutPlan{
		Name:             name,
		Description:      description,
		TrainerID:        trainerID,
		WorkoutExercises: []models.WorkoutExercise{},
	}

	// Map inputs to the WorkoutExercise model.
	for i, ex := range exercises {
		link := models.WorkoutExercise{
			ExerciseID: ex.ExerciseID,
			Sets:       ex.Sets,
			Reps:       ex.Reps,
			Notes:      ex.Notes,
			Order:      i + 1, // Automatically maintain sequence order.
		}

		plan.WorkoutExercises = append(plan.WorkoutExercises, link)
	}

	// Persist the entire plan graph.
	return s.Repo.Create(&plan)
}

// GetFullDetails retrieves a workout plan by its ID with all exercises preloaded.
func (s *WorkoutService) GetFullDetails(id uint) (*models.WorkoutPlan, error) {
	return s.Repo.GetByID(id)
}

// ListAll retrieves all workout plans available in the system.
func (s *WorkoutService) ListAll() ([]models.WorkoutPlan, error) {
	return s.Repo.GetAll()
}

// UpdatePlan modifies an existing workout plan and replaces its exercise list.
func (s *WorkoutService) UpdatePlan(planID uint, name, description string, exIDs, sets, reps, rowNotes []string) error {
	// Retrieve the existing plan.
	plan, err := s.Repo.GetByID(planID)
	if err != nil {
		return err
	}

	// Update metadata.
	plan.Name = name
	plan.Description = description

	// Remove old exercise associations.
	if err := s.Repo.ClearExercises(planID); err != nil {
		return err
	}

	// Rebuild the exercise list from the provided form data.
	var exercises []models.WorkoutExercise
	for i, exIDStr := range exIDs {
		exID, _ := strconv.ParseUint(exIDStr, 10, 32)
		sVal, _ := strconv.Atoi(sets[i])
		rVal, _ := strconv.Atoi(reps[i])

		nVal := ""
		if i < len(rowNotes) {
			nVal = rowNotes[i]
		}

		exercises = append(exercises, models.WorkoutExercise{
			WorkoutPlanID: planID,
			ExerciseID:    uint(exID),
			Sets:          sVal,
			Reps:          rVal,
			Notes:         nVal,
			Order:         i + 1,
		})
	}

	plan.WorkoutExercises = exercises
	return s.Repo.Update(plan)
}

// DeletePlan removes a workout plan from the system.
func (s *WorkoutService) DeletePlan(id uint) error {
	return s.Repo.Delete(id)
}
