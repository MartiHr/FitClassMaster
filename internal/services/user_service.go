package services

import (
	"FitClassMaster/internal/models"
	"FitClassMaster/internal/repositories"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

// UserService provides administrative and profile-related operations for users.
type UserService struct {
	Repo *repositories.UserRepo
}

// NewUserService creates a new instance of UserService.
func NewUserService(repo *repositories.UserRepo) *UserService {
	return &UserService{Repo: repo}
}

// GetProfile retrieves a user's profile information by ID.
func (s *UserService) GetProfile(id uint) (*models.User, error) {
	return s.Repo.GetById(id)
}

// UpdateProfile updates the name information for a user.
func (s *UserService) UpdateProfile(id uint, firstName, lastName string) error {
	return s.Repo.UpdateInfo(id, firstName, lastName)
}

// ChangePassword updates a user's password after verifying the current one.
func (s *UserService) ChangePassword(userID uint, currentPwd, newPwd string) error {
	user, err := s.Repo.GetById(userID)
	if err != nil {
		return err
	}

	// Verify the current password.
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(currentPwd)); err != nil {
		return errors.New("Current password is incorrect")
	}

	// Generate hash for the new password.
	hash, err := bcrypt.GenerateFromPassword([]byte(newPwd), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return s.Repo.UpdatePassword(userID, string(hash))
}

// ListAllUsers retrieves all users in the system.
func (s *UserService) ListAllUsers() ([]models.User, error) {
	return s.Repo.GetAll()
}

// PromoteUser upgrades a user's role to 'trainer'.
func (s *UserService) PromoteUser(userID uint) error {
	return s.Repo.UpdateRole(userID, models.RoleTrainer)
}

// DemoteUser sets a user's role back to 'member'.
func (s *UserService) DemoteUser(userID uint) error {
	return s.Repo.UpdateRole(userID, models.RoleMember)
}

// DeleteUser permanently removes a user from the system.
func (s *UserService) DeleteUser(userID uint) error {
	return s.Repo.Delete(userID)
}
