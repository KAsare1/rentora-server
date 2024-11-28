package repository

import (
	"errors"
	"fmt"
	"strings"

	"rentora-go/internal/model"

	"gorm.io/gorm"
)

type UserRepository interface {
	GetUserByEmail(email string) (*model.User, error)
	CreateUser(user *model.User) error
	UpdateUser(user *model.User) error
	ListUsers(limit, offset int) ([]model.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetUserByEmail(email string) (*model.User, error) {
	var user model.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user with email %s: %w", email, err)
	}
	return &user, nil
}

func (r *userRepository) CreateUser(user *model.User) error {
    if err := r.db.Create(user).Error; err != nil {
        if strings.Contains(err.Error(), "duplicate key value") || strings.Contains(err.Error(), "unique constraint") {
            return errors.New("email already registered")
        }
        return err
    }
    return nil
}

func (r *userRepository) UpdateUser(user *model.User) error {
	if err := r.db.Save(user).Error; err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}

func (r *userRepository) ListUsers(limit, offset int) ([]model.User, error) {
	var users []model.User
	err := r.db.Limit(limit).Offset(offset).Find(&users).Error
	if err != nil {
		return nil, fmt.Errorf("failed to fetch users: %w", err)
	}
	return users, nil
}
