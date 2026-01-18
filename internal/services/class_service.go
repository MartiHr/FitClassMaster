package services

import (
	"FitClassMaster/internal/models"
	"FitClassMaster/internal/repositories"
)

type ClassService struct {
	Repo *repositories.ClassRepo
}

func NewClassService(repo *repositories.ClassRepo) *ClassService {
	return &ClassService{Repo: repo}
}

func (s *ClassService) GetAvailableClasses() ([]models.Class, error) {
	return s.Repo.GetAll()
}
