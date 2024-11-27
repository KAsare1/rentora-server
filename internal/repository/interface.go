package repository

import "rentora-go/internal/model"

type UserRepository interface {
	Create(user *model.User) error
	FindByEmail(email string) (*model.User, error)
	FindById(id uint) (*model.User, error)
	Update(user *model.User) error
	Delete(id uint) error
}

