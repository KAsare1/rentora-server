package handler

import (
	"encoding/json"
	"net/http"
	"time"
	"log"

	"rentora-go/internal/model"
	"rentora-go/internal/service"

)

type AuthHandler struct {
	authService service.AuthService
	
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Login handles user login requests
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req model.LoginRequest

	// Decode the request payload into the LoginRequest struct
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Authenticate the user
	accessToken, refreshToken, err := h.authService.Authenticate(req.Email, req.Password)
	if err != nil {
		http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
		return
	}

	// Respond with the tokens
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req model.User

	// Decode the request payload
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Failed to decode request: %v", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.FirstName == "" || req.LastName == "" || req.Email == "" || req.Password == "" || req.PhoneNumber == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	// Check if email already exists
	existingUser, _ := h.authService.GetUserByEmail(req.Email)
	if existingUser != nil {
		http.Error(w, "Email already in use", http.StatusConflict)
		return
	}

	// Hash the user's password
	if err := req.HashPassword(); err != nil {
		http.Error(w, "Error securing password", http.StatusInternalServerError)
		return
	}

	// Set default values
	req.Role = "customer" // Default role
	req.RegistrationDate = time.Now()
	req.IsVerified = false // Initially unverified
	req.IsActive = true

	// Save the user via the AuthService
	if err := h.authService.CreateUser(&req); err != nil {
		http.Error(w, "Failed to create user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with the newly created user's public details
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(req.ToDTO())
}

// Refresh handles token refresh requests
func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req model.RefreshTokenRequest

	// Decode the request payload into the RefreshTokenRequest struct
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Refresh the tokens
	accessToken, refreshToken, err := h.authService.RefreshToken(req.RefreshToken)
	if err != nil {
		http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
		return
	}

	// Respond with the new tokens
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}
