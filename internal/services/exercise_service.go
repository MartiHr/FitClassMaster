package services

import (
	"FitClassMaster/internal/models"
	"FitClassMaster/internal/repositories"
)

// ExerciseService manages the catalog of available exercises.
type ExerciseService struct {
	Repo *repositories.ExerciseRepo
}

// NewExerciseService creates a new instance of ExerciseService.
func NewExerciseService(repo *repositories.ExerciseRepo) *ExerciseService {
	return &ExerciseService{Repo: repo}
}

// GetAll retrieves all exercises from the repository.
func (s *ExerciseService) GetAll() ([]models.Exercise, error) {
	return s.Repo.GetAll()
}

// Create adds a new exercise to the catalog.
func (s *ExerciseService) Create(name, description, muscleGroup, equipment, videoURL string) error {
	exercise := &models.Exercise{
		Name:        name,
		Description: description,
		MuscleGroup: muscleGroup,
		Equipment:   equipment,
		VideoURL:    videoURL,
	}
	return s.Repo.Create(exercise)
}

// GetByID retrieves a single exercise by its ID.
func (s *ExerciseService) GetByID(id uint) (*models.Exercise, error) {
	return s.Repo.GetByID(id)
}

// Update modifies an existing exercise's details.
func (s *ExerciseService) Update(id uint, name, description, muscleGroup, equipment, videoURL string) error {
	exercise, err := s.Repo.GetByID(id)
	if err != nil {
		return err
	}

	exercise.Name = name
	exercise.Description = description
	exercise.MuscleGroup = muscleGroup
	exercise.Equipment = equipment
	exercise.VideoURL = videoURL

	return s.Repo.Update(exercise)
}

// Delete removes an exercise from the catalog.
func (s *ExerciseService) Delete(id uint) error {
	return s.Repo.Delete(id)
}
