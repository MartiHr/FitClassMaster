package services

import (
	"FitClassMaster/internal/models"
	"FitClassMaster/internal/repositories"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	Repo *repositories.UserRepo
}

func NewUserService(repo *repositories.UserRepo) *UserService {
	return &UserService{Repo: repo}
}

func (s *UserService) GetProfile(id uint) (*models.User, error) {
	return s.Repo.GetById(id)
}

func (s *UserService) UpdateProfile(id uint, firstName, lastName string) error {
	return s.Repo.UpdateInfo(id, firstName, lastName)
}

func (s *UserService) ChangePassword(userID uint, currentPwd, newPwd string) error {
	user, err := s.Repo.GetById(userID)
	if err != nil {
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(currentPwd)); err != nil {
		return errors.New("Current password is incorrect")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(newPwd), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return s.Repo.UpdatePassword(userID, string(hash))
}

func (s *UserService) ListAllUsers() ([]models.User, error) {
	return s.Repo.GetAll()
}

func (s *UserService) PromoteUser(userID uint) error {
	return s.Repo.UpdateRole(userID, models.RoleTrainer)
}

func (s *UserService) DemoteUser(userID uint) error {
	return s.Repo.UpdateRole(userID, models.RoleMember)
}

func (s *UserService) DeleteUser(userID uint) error {
	return s.Repo.Delete(userID)
}
