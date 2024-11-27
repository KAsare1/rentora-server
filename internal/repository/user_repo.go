package repository

import (
	"errors"
	"gorm.io/gorm"
	"rentora-go/internal/model"
)


type UserRepositoryImpl struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepositoryImpl {
	return &UserRepositoryImpl{db}
}



func (r *UserRepositoryImpl) Create(user *model.User) error {
	if error := user.HashPassword(); error != nil {
		return error
	}

	var existingUser model.User
	if result := r.db.Where("email = ?", user.Email).First(&existingUser); result.Error == nil {
		return errors.New("user with email already exists")
	}

	return r.db.Create(user).Error
}



func (r *UserRepositoryImpl) FindByEmail(email string) (*model.User, error) {
	var user model.User
	if result := r.db.Where("email = ?", email).First(&user); result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}


func (r *UserRepositoryImpl) FindById(id uint) (*model.User, error) {
	var user model.User
	if result := r.db.Where("id = ?", id).First(&user); result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}




