package handler

import (
	"encoding/json"
	"net/http"

	"rentora-go/internal/model"
	"rentora-go/internal/repository"
	"rentora-go/internal/service"
)

type AuthHandler struct {
	userRepo   repository.UserRepository
	jwtService service.JWTService
}

func NewAuthHandler(userRepo repository.UserRepository, jwtService service.JWTService) *AuthHandler {
	return &AuthHandler{userRepo: userRepo, jwtService: jwtService}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var creds model.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	user, err := h.userRepo.GetUserByEmail(creds.Email)
	if err != nil || !user.CheckPassword(creds.Password) {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	accessToken, refreshToken, _ := h.jwtService.GenerateTokens(user.Email)
	response := map[string]string{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	}

	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var request model.RefreshTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	email, err := h.jwtService.ValidateRefreshToken(request.RefreshToken)
	if err != nil {
		http.Error(w, "Invalid refresh token", http.StatusUnauthorized)
		return
	}

	newAccessToken, _, _ := h.jwtService.GenerateTokens(email)
	response := map[string]string{
		"access_token": newAccessToken,
	}

	json.NewEncoder(w).Encode(response)
}
