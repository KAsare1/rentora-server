package service

import "rentora-go/internal/model"


type AuthService interface {
	Register(user *model.User) (*model.UserDTO, error)
	Login(email, password string) (*model.UserDTO, error)
	GenerateToken(user *model.UserDTO) (string, error)
}