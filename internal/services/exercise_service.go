package services

import (
	"FitClassMaster/internal/models"
	"FitClassMaster/internal/repositories"
)

type ExerciseService struct {
	Repo *repositories.ExerciseRepo
}

func NewExerciseService(repo *repositories.ExerciseRepo) *ExerciseService {
	return &ExerciseService{Repo: repo}
}

func (s *ExerciseService) GetAll() ([]models.Exercise, error) {
	return s.Repo.GetAll()
}

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

func (s *ExerciseService) GetByID(id uint) (*models.Exercise, error) {
	return s.Repo.GetByID(id)
}

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

func (s *ExerciseService) Delete(id uint) error {
	return s.Repo.Delete(id)
}
