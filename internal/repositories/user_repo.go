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

func (r *UserRepo) GetById(id int) (user *models.User, err error) {
	var u models.User

	if err := config.DB.First(&u, id).Error; err != nil {
		return nil, err
	}

	return &u, nil
}

func (r *UserRepo) GetByEmail(email string) (user *models.User, err error) {
	var u models.User

	if err := config.DB.Where("email = ?", email).First(&user).Error; err != nil {
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
