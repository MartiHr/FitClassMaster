package services

import (
	"FitClassMaster/internal/models"
	"FitClassMaster/internal/repositories"

	"errors"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo *repositories.UserRepo
}

func NewAuthService(r *repositories.UserRepo) *AuthService {
	return &AuthService{repo: r}
}

func (as *AuthService) Register(u *models.User, plainPassword string) error {
	// check whether it exists
	if _, err := as.repo.GetByEmail(u.Email); err == nil {
		return errors.New("email already used")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	u.Password = string(hash)

	return as.repo.Create(u)
}

func (as *AuthService) Login(email, password string) (*models.User, error) {
	user, err := as.repo.GetByEmail(email)

	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, err
	}

	return user, nil
}
