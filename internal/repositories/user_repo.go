// Package repositories provides the data access layer for interacting with the database.
package repositories

import (
	"FitClassMaster/internal/config"
	"FitClassMaster/internal/models"
)

// UserRepo handles database operations for the User model.
type UserRepo struct{}

// NewUserRepo creates a new instance of UserRepo.
func NewUserRepo() *UserRepo {
	return &UserRepo{}
}

// Create inserts a new user into the database.
func (r *UserRepo) Create(u *models.User) error {
	return config.DB.Create(u).Error
}

// GetById retrieves a user by their primary key ID.
func (r *UserRepo) GetById(id uint) (user *models.User, err error) {
	var u models.User
	if err := config.DB.First(&u, id).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

// GetByEmail retrieves a user by their unique email address.
func (r *UserRepo) GetByEmail(email string) (user *models.User, err error) {
	var u models.User
	if err := config.DB.Where("email = ?", email).First(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

// UpdateInfo updates the first and last name of a user.
func (r *UserRepo) UpdateInfo(id uint, firstName, lastName string) error {
	return config.DB.Model(&models.User{}).Where("id = ?", id).Updates(map[string]interface{}{
		"first_name": firstName,
		"last_name":  lastName,
	}).Error
}

// UpdatePassword updates the hashed password for a user.
func (r *UserRepo) UpdatePassword(id uint, hashedPwd string) error {
	return config.DB.Model(&models.User{}).Where("id = ?", id).Update("password", hashedPwd).Error
}

// GetAll fetches all users in the system, ordered by their creation date (newest first).
func (r *UserRepo) GetAll() ([]models.User, error) {
	var users []models.User
	err := config.DB.Order("created_at desc").Find(&users).Error
	return users, err
}

// UpdateRole updates a user's authorization role.
func (r *UserRepo) UpdateRole(userID uint, newRole models.Role) error {
	return config.DB.Model(&models.User{}).Where("id = ?", userID).Update("role", newRole).Error
}

// Delete removes a user from the database permanently.
func (r *UserRepo) Delete(userID uint) error {
	return config.DB.Unscoped().Delete(&models.User{}, userID).Error
}
