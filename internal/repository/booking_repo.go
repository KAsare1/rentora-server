package repository

import (
	"rentora-go/internal/model"
	"gorm.io/gorm"
)

type BookingRepository interface {
	CreateBooking(booking *model.Booking) error
	GetBookingsByUserID(userID uint) ([]model.Booking, error)
	GetBookingByID(bookingID uint) (*model.Booking, error)
	UpdateBooking(booking *model.Booking) error
	DeleteBooking(bookingID uint) error
}

type bookingRepository struct {
	db *gorm.DB
}

func NewBookingRepository(db *gorm.DB) BookingRepository {
	return &bookingRepository{db: db}
}

func (r *bookingRepository) CreateBooking(booking *model.Booking) error {
	if err := r.db.Create(booking).Error; err != nil {
		return err
	}
	return nil
}

func (r *bookingRepository) GetBookingsByUserID(userID uint) ([]model.Booking, error) {
	var bookings []model.Booking
	if err := r.db.Where("user_id = ?", userID).Find(&bookings).Error; err != nil {
		return nil, err
	}
	return bookings, nil
}

func (r *bookingRepository) GetBookingByID(bookingID uint) (*model.Booking, error) {
	var booking model.Booking
	if err := r.db.First(&booking, bookingID).Error; err != nil {
		return nil, err
	}
	return &booking, nil
}

func (r *bookingRepository) UpdateBooking(booking *model.Booking) error {
	if err := r.db.Save(booking).Error; err != nil {
		return err
	}
	return nil
}

func (r *bookingRepository) DeleteBooking(bookingID uint) error {
	if err := r.db.Delete(&model.Booking{}, bookingID).Error; err != nil {
		return err
	}
	return nil
}
