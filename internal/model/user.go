package model

import (
	"golang.org/x/crypto/bcrypt"
	"time"
	"log"
)

// User represents a user in the system
type User struct {
	ID                    uint           `json:"id"`
	FirstName             string         `json:"first_name"`
	LastName              string         `json:"last_name"`
	OtherName             string         `json:"other_name"`
	Email                 string         `json:"email"`
	Password              string         `json:"password"`
	Role                  string         `json:"role"`

	// Personal Information
	PhoneNumber           string         `gorm:"not null" json:"phone_number"`
	DateOfBirth           *Date      `gorm:"type:date" json:"date_of_birth"` // Changed to Date for proper date handling
	Address               string         `json:"address"`
	City                  string         `json:"city"`
	Region                string         `json:"region"`
	Country               string         `json:"country"`
	PostalCode            string         `json:"postal_code"`

	// Driver's License Information
	DriversLicenseNumber  string         `json:"drivers_license_number"`
	DriversLicenseState   string         `json:"drivers_license_state"`
	DriversLicenseExpiration Date   `gorm:"type:date" json:"drivers_license_expiration"` // Changed to Date

	// Account Status and Verification
	IsVerified            bool           `gorm:"default:false" json:"is_verified"`
	IsActive              bool           `gorm:"default:true" json:"is_active"`
	RegistrationDate      time.Time      `gorm:"autoCreateTime" json:"registration_date"`
	LastLoginDate         *time.Time      `json:"last_login_date"`

	// Payment and Billing
	PaymentMethod         string         `json:"payment_method"`
	HasOutstandingBalance bool           `gorm:"default:false" json:"has_outstanding_balance"`
	AccountCredit         float64        `gorm:"default:0" json:"account_credit"`

	// Rental History and Preferences
	TotalRentals          int            `gorm:"default:0" json:"total_rentals"`
	CurrentRentalCount    int            `gorm:"default:0" json:"current_rental_count"`
	PreferredVehicleType  string         `json:"preferred_vehicle_type"`

	// Security and Compliance
	AcceptedTermsOfService bool          `gorm:"default:false" json:"accepted_terms_of_service"`
	TermsAcceptedDate      *time.Time     `json:"terms_accepted_date"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// HashPassword hashes the user's password using bcrypt.
func (u *User) HashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

func (u *User) CheckPassword(password string) bool {
    log.Println("Provided password:", password)
    log.Println("Hashed password in DB:", u.Password)
    err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
    if err != nil {
        log.Println("Password comparison failed:", err)
    }
    return err == nil
}
// UserDTO is a data transfer object for the User model (to return to the API).
type UserDTO struct {
	ID        uint   `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	OtherName string `json:"other_name"`
	Email     string `json:"email"`
	Role      string `json:"role"`
}

// ToDTO converts a User model into a UserDTO.
func (u *User) ToDTO() UserDTO {
	return UserDTO{
		ID:        u.ID,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		OtherName: u.OtherName,
		Email:     u.Email,
		Role:      u.Role,
	}
}


type UserUpdateRequest struct {
	Email                 string    `json:"email" validate:"omitempty,email"`
	Password              string    `json:"password" validate:"omitempty,min=8"`
	PhoneNumber           string    `json:"phone_number" validate:"omitempty,e164"`
	Address               string    `json:"address" validate:"omitempty"`
	City                  string    `json:"city" validate:"omitempty"`
	Region                string    `json:"region" validate:"omitempty"`
	Country               string    `json:"country" validate:"omitempty,iso3166_1_alpha2"`
	PostalCode            string    `json:"postal_code" validate:"omitempty"`
	DriversLicenseExpiration *Date  `json:"drivers_license_expiration" validate:"omitempty"`
	PaymentMethod         string    `json:"payment_method" validate:"omitempty,oneof=credit_card debit_card paypal bank_transfer"`
	PreferredVehicleType  string    `json:"preferred_vehicle_type" validate:"omitempty,oneof=sedan suv truck compact luxury electric hybrid"`
}

// LoginRequest represents the request payload for the login endpoint
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// RefreshTokenRequest represents the request payload for the refresh token endpoint
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}