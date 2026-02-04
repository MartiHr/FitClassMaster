// Package services contains the business logic layer of the application.
package services

import (
	"FitClassMaster/internal/models"
	"FitClassMaster/internal/repositories"

	"errors"

	"golang.org/x/crypto/bcrypt"
)

// AuthService handles user registration and authentication logic.
type AuthService struct {
	repo *repositories.UserRepo
}

// NewAuthService creates a new instance of AuthService.
func NewAuthService(r *repositories.UserRepo) *AuthService {
	return &AuthService{repo: r}
}

// Register hashes the user's password and saves the new user to the database.
// It returns an error if the email is already in use.
func (as *AuthService) Register(u *models.User, plainPassword string) error {
	// Check if a user with the same email already exists.
	if _, err := as.repo.GetByEmail(u.Email); err == nil {
		return errors.New("email already used")
	}

	// Hash the plain text password.
	hash, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.Password = string(hash)

	// Persist the user record.
	return as.repo.Create(u)
}

// Login validates user credentials against the database.
// It returns the user object if successful, or an error if validation fails.
func (as *AuthService) Login(email, password string) (*models.User, error) {
	// Retrieve the user by email.
	user, err := as.repo.GetByEmail(email)
	if err != nil {
		return nil, err
	}

	// Compare the provided password with the stored hash.
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, err
	}

	return user, nil
}
