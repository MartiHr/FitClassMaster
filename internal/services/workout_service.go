package services

import (
	"FitClassMaster/internal/models"
	"FitClassMaster/internal/repositories"
	"strconv"
)

type WorkoutService struct {
	Repo *repositories.WorkoutRepo
}

func NewWorkoutService(repo *repositories.WorkoutRepo) *WorkoutService {
	return &WorkoutService{Repo: repo}
}

// ExerciseInput represents one row in the "Add Workout" form
type ExerciseInput struct {
	ExerciseID uint
	Sets       int
	Reps       int
	Notes      string
}

// CreatePlan constructs the full object graph and saves it
func (s *WorkoutService) CreatePlan(name, description string, trainerID uint, exercises []ExerciseInput) error {

	// Create the Parent Plan
	plan := models.WorkoutPlan{
		Name:             name,
		Description:      description,
		TrainerID:        trainerID,
		WorkoutExercises: []models.WorkoutExercise{}, // Initialize slice
	}

	// Loop through inputs and create the Child objects
	for i, ex := range exercises {
		link := models.WorkoutExercise{
			ExerciseID: ex.ExerciseID,
			Sets:       ex.Sets,
			Reps:       ex.Reps,
			Notes:      ex.Notes,
			Order:      i + 1, // Automatically set order based on list position
		}

		plan.WorkoutExercises = append(plan.WorkoutExercises, link)
	}

	// Save everything at once using the Repo
	return s.Repo.Create(&plan)
}

func (s *WorkoutService) GetFullDetails(id uint) (*models.WorkoutPlan, error) {
	return s.Repo.GetByID(id)
}

func (s *WorkoutService) ListAll() ([]models.WorkoutPlan, error) {
	return s.Repo.GetAll()
}

// UpdatePlan updates the plan metadata and replaces all exercises
func (s *WorkoutService) UpdatePlan(planID uint, name, description string, exIDs, sets, reps, rowNotes []string) error {
	// Get the plan
	plan, err := s.Repo.GetByID(planID)
	if err != nil {
		return err
	}

	// Update Basic Info
	plan.Name = name
	plan.Description = description // Matches your Struct field

	// Clear existing exercises
	if err := s.Repo.ClearExercises(planID); err != nil {
		return err
	}

	// Re-add Exercises from Form Data
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
			Notes:         nVal, // This belongs to WorkoutExercise
			Order:         i + 1,
		})
	}

	plan.WorkoutExercises = exercises
	return s.Repo.Update(plan)
}

func (s *WorkoutService) DeletePlan(id uint) error {
	return s.Repo.Delete(id)
}
