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
	// Change r.DB to config.DB
	return config.DB.Model(&models.User{}).Where("id = ?", id).Update("password", hashedPwd).Error
}
