package service

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTService interface {
	GenerateTokens(email string) (string, string, error)
	ValidateRefreshToken(refreshToken string) (string, error)
}

type jwtService struct {
	secretKey       string
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

func NewJWTService(secretKey string, accessTokenTTL, refreshTokenTTL int) JWTService {
	return &jwtService{
		secretKey:       secretKey,
		accessTokenTTL:  time.Duration(accessTokenTTL) * time.Hour,
		refreshTokenTTL: time.Duration(refreshTokenTTL) * time.Hour,
	}
}

func (j *jwtService) GenerateTokens(email string) (string, string, error) {
	now := time.Now()

	// Access Token
	accessTokenClaims := jwt.MapClaims{
		"sub": email,
		"exp": now.Add(j.accessTokenTTL).Unix(),
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
	accessTokenString, err := accessToken.SignedString([]byte(j.secretKey))
	if err != nil {
		return "", "", err
	}

	// Refresh Token
	refreshTokenClaims := jwt.MapClaims{
		"sub": email,
		"exp": now.Add(j.refreshTokenTTL).Unix(),
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(j.secretKey))
	if err != nil {
		return "", "", err
	}

	return accessTokenString, refreshTokenString, nil
}

func (j *jwtService) ValidateRefreshToken(refreshToken string) (string, error) {
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.secretKey), nil
	})

	if err != nil || !token.Valid {
		return "", errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("invalid claims")
	}

	email, ok := claims["sub"].(string)
	if !ok {
		return "", errors.New("invalid subject claim")
	}

	return email, nil
}
