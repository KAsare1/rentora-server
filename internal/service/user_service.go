package service

import (
	"errors"
	"time"

	"log"

	"rentora-go/internal/model"
	"rentora-go/internal/repository"

	"github.com/golang-jwt/jwt/v5"
)

type AuthService interface {
	Authenticate(email, password string) (string, string, error)
	RefreshToken(refreshToken string) (string, string, error)
	CreateUser(user *model.User) error                // Add this
	GetUserByEmail(email string) (*model.User, error)
}

type authService struct {
	userRepo repository.UserRepository
	jwtKey   []byte
}

func NewAuthService(userRepo repository.UserRepository, jwtKey []byte) AuthService {
	return &authService{userRepo: userRepo, jwtKey: jwtKey}
}

type Claims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

func (s *authService) Authenticate(email, password string) (string, string, error) {
    user, err := s.userRepo.GetUserByEmail(email)
    if err != nil {
        log.Println("Error fetching user by email:", err)
        return "", "", errors.New("invalid email or password")
    }
    log.Println("User fetched:", user)

    if !user.CheckPassword(password) {
        log.Println("Password check failed")
        return "", "", errors.New("invalid email or password")
    }
    log.Println("Password check passed")

    accessToken, err := s.generateToken(user.ID, email, time.Minute*15)
    if err != nil {
        log.Println("Error generating access token:", err)
        return "", "", err
    }

    refreshToken, err := s.generateToken(user.ID, email, time.Hour*24*7)
    if err != nil {
        log.Println("Error generating refresh token:", err)
        return "", "", err
    }

    return accessToken, refreshToken, nil
}

func (s *authService) RefreshToken(refreshToken string) (string, string, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(refreshToken, claims, func(token *jwt.Token) (interface{}, error) {
		return s.jwtKey, nil
	})

	if err != nil || !token.Valid {
		return "", "", errors.New("invalid refresh token")
	}

	newAccessToken, err := s.generateToken(claims.UserID, claims.Email, time.Minute*15)
	if err != nil {
		return "", "", err
	}

	newRefreshToken, err := s.generateToken(claims.UserID, claims.Email, time.Hour*24*7)
	if err != nil {
		return "", "", err
	}

	return newAccessToken, newRefreshToken, nil
}

func (s *authService) generateToken(userID uint, email string, duration time.Duration) (string, error) {
	claims := &Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtKey)
}


func (s *authService) GetUserByEmail(email string) (*model.User, error) {
	return s.userRepo.GetUserByEmail(email)
}

func (s *authService) CreateUser(user *model.User) error {
	return s.userRepo.CreateUser(user)
}
