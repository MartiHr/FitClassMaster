package repositories

import (
	"FitClassMaster/internal/config"
	"FitClassMaster/internal/models"
)

type UserRepo struct{}

func NewUserRepo() *UserRepo {
	return &UserRepo{}
}

func (r *UserRepo) Create(u *models.User) error {
	return config.DB.Create(u).Error
}

func (r *UserRepo) GetById(id uint) (user *models.User, err error) {
	var u models.User

	if err := config.DB.First(&u, id).Error; err != nil {
		return nil, err
	}

	return &u, nil
}

func (r *UserRepo) GetByEmail(email string) (user *models.User, err error) {
	var u models.User

	if err := config.DB.Where("email = ?", email).First(&u).Error; err != nil {
		return nil, err
	}

	return &u, nil
}

func (r *UserRepo) GetByName(name string) (user *models.User, err error) {
	var u models.User

	if err := config.DB.Where("name = ?", name).First(&u).Error; err != nil {
		return nil, err
	}

	return &u, nil
}

// UpdateInfo updates the first and last name of a user
func (r *UserRepo) UpdateInfo(id uint, firstName, lastName string) error {
	// Change r.DB to config.DB
	return config.DB.Model(&models.User{}).Where("id = ?", id).Updates(map[string]interface{}{
		"first_name": firstName,
		"last_name":  lastName,
	}).Error
}

// UpdatePassword updates the hashed password in the database
func (r *UserRepo) UpdatePassword(id uint, hashedPwd string) error {
	return config.DB.Model(&models.User{}).Where("id = ?", id).Update("password", hashedPwd).Error
}

// GetAll fetches every user in the system, ordered by newest first
func (r *UserRepo) GetAll() ([]models.User, error) {
	var users []models.User
	err := config.DB.Order("created_at desc").Find(&users).Error
	return users, err
}

// UpdateRole changes a user's role (e.g., "member" -> "trainer")
func (r *UserRepo) UpdateRole(userID uint, newRole models.Role) error {
	return config.DB.Model(&models.User{}).Where("id = ?", userID).Update("role", newRole).Error
}

// Delete removes a user permanently
func (r *UserRepo) Delete(userID uint) error {
	return config.DB.Unscoped().Delete(&models.User{}, userID).Error
}
