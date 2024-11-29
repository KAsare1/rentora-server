package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"rentora-go/internal/model"
	"rentora-go/internal/service"
)

type AuthHandler struct {
	authService service.AuthService
	
}

type UpdateUserRequest struct {
    Email                   *string `json:"email,omitempty"`
    Password                *string `json:"password,omitempty"`
    PhoneNumber             *string `json:"phone_number,omitempty"`
    Address                 *string `json:"address,omitempty"`
    City                    *string `json:"city,omitempty"`
    Region                  *string `json:"region,omitempty"`
    Country                 *string `json:"country,omitempty"`
    PostalCode              *string `json:"postal_code,omitempty"`
    DriversLicenseExpiration *string `json:"drivers_license_expiration,omitempty"`
    PaymentMethod           *string `json:"payment_method,omitempty"`
    PreferredVehicleType    *string `json:"preferred_vehicle_type,omitempty"`
}



func (req *UpdateUserRequest) ToServiceUpdateUserRequest() service.UpdateUserRequest {
    return service.UpdateUserRequest{
        Email:                   req.Email,
        Password:                req.Password,
        PhoneNumber:             req.PhoneNumber,
        Address:                 req.Address,
        City:                    req.City,
        Region:                  req.Region,
        Country:                 req.Country,
        PostalCode:              req.PostalCode,
        DriversLicenseExpiration: req.DriversLicenseExpiration,
        PaymentMethod:           req.PaymentMethod,
        PreferredVehicleType:    req.PreferredVehicleType,
    }
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




func (h *AuthHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
    userID, ok := r.Context().Value("user_id").(uint)
    if !ok {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    var updateReq UpdateUserRequest
    if err := json.NewDecoder(r.Body).Decode(&updateReq); err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    // Convert to service.UpdateUserRequest
    serviceUpdateReq := updateReq.ToServiceUpdateUserRequest()

    // Call the service method with the converted request
    err := h.authService.UpdateUser(userID, serviceUpdateReq)
    if err != nil {
        switch {
        case err.Error() == "user not found":
            http.Error(w, "User not found", http.StatusNotFound)
        case err.Error() == "email already in use":
            http.Error(w, "Email already in use", http.StatusConflict)
        case strings.Contains(err.Error(), "invalid"):
            http.Error(w, err.Error(), http.StatusBadRequest)
        default:
            http.Error(w, "Failed to update user", http.StatusInternalServerError)
        }
        return
    }

    // Fetch the updated user to return
    updatedUser, err := h.authService.GetUserByID(userID)
    if err != nil {
        http.Error(w, "Failed to retrieve updated user", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(updatedUser.ToDTO())
}