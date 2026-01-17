package services

import (
	"FitClassMaster/internal/models"
	"FitClassMaster/internal/repositories"
)

type UserService struct {
	Repo *repositories.UserRepo
}

func NewUserService(repo *repositories.UserRepo) *UserService {
	return &UserService{Repo: repo}
}

func (s *UserService) GetProfile(id uint) (*models.User, error) {
	// You could add logic here later, like checking if account is active
	return s.Repo.GetById(id)
}
