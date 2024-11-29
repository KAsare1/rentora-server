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

	UpdateUser(userID uint, updateReq UpdateUserRequest) error
	GetUserByID(userID uint) (*model.User, error)
}

type UpdateUserRequest struct {
	Email                    *string
	Password                 *string
	PhoneNumber              *string
	Address                  *string
	City                     *string
	Region                   *string
	Country                  *string
	PostalCode               *string
	DriversLicenseExpiration *string
	PaymentMethod            *string
	PreferredVehicleType     *string
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



func (s *authService) UpdateUser(userID uint, updateReq UpdateUserRequest) error {
    // Fetch the existing user
    user, err := s.userRepo.GetUserByID(userID)
    if err != nil {
        return errors.New("user not found")
    }

    // Validate and update email (check for uniqueness)
    if updateReq.Email != nil {
        // Check if email is already in use by another user
        existingUser, _ := s.userRepo.GetUserByEmail(*updateReq.Email)
        if existingUser != nil && existingUser.ID != userID {
            return errors.New("email already in use")
        }
        user.Email = *updateReq.Email
    }

    // Update password with hashing
    if updateReq.Password != nil {
        user.Password = *updateReq.Password
        if err := user.HashPassword(); err != nil {
            return errors.New("failed to update password")
        }
    }

    // Update other fields
    if updateReq.PhoneNumber != nil {
        user.PhoneNumber = *updateReq.PhoneNumber
    }
    if updateReq.Address != nil {
        user.Address = *updateReq.Address
    }
    if updateReq.City != nil {
        user.City = *updateReq.City
    }
    if updateReq.Region != nil {
        user.Region = *updateReq.Region
    }
    if updateReq.Country != nil {
        user.Country = *updateReq.Country
    }
    if updateReq.PostalCode != nil {
        user.PostalCode = *updateReq.PostalCode
    }

    // Handle DriversLicenseExpiration
    if updateReq.DriversLicenseExpiration != nil {
        parsedTime, err := time.Parse(time.DateOnly, *updateReq.DriversLicenseExpiration)
        if err != nil {
            return errors.New("invalid drivers license expiration date format")
        }
        expirationDate := model.Date{Time: parsedTime}
        user.DriversLicenseExpiration = expirationDate
    }

    if updateReq.PaymentMethod != nil {
        user.PaymentMethod = *updateReq.PaymentMethod
    }
    if updateReq.PreferredVehicleType != nil {
        user.PreferredVehicleType = *updateReq.PreferredVehicleType
    }

    // Save the updated user
    return s.userRepo.UpdateUser(user)
}

func (s *authService) GetUserByID(userID uint) (*model.User, error) {
    return s.userRepo.GetUserByID(userID)
}