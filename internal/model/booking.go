package model

import (
	"time"
)

type Booking struct {
	ID            uint      `json:"id"`
	UserID        uint      `json:"user_id"`
	CarID         uint      `json:"car_id"`
	StartDate     time.Time `json:"start_date"`
	EndDate       time.Time `json:"end_date"`
	TotalAmount   float64   `json:"total_amount"`
	Status        string    `json:"status"` // e.g., "Pending", "Accepted", "Declined", "Completed", "Cancelled"
	PaymentMethod string    `json:"payment_method"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
