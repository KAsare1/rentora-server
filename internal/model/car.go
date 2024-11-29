package model

import "time"

type Car struct {
	ID           uint      `json:"id"`
	OwnerID      uint      `json:"owner_id"`
	Make         string    `json:"make"`
	Model        string    `json:"model"`
	Year         int       `json:"year"`
	PricePerDay  float64   `json:"price_per_day"`
	Availability bool      `gorm:"default:true" json:"availability"`
	Location     string    `json:"location"`
	Description  string    `json:"description"`
	ImageURL     string    `json:"image_url"` // Optional: Add image URLs for car photos
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	Owner User `gorm:"foreignKey:OwnerID" json:"owner"`
}
