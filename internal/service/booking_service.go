package service

import (
	"errors"
	"rentora-go/internal/model"
	"rentora-go/internal/repository"
)

type BookingService struct {
	repo repository.BookingRepository
}

func NewBookingService(repo repository.BookingRepository) *BookingService {
	return &BookingService{repo: repo}
}

func (s *BookingService) CreateBooking(booking *model.Booking) error {
	// Add additional logic to check car availability
	booking.Status = "Pending" // Default status when booking is created
	return s.repo.CreateBooking(booking)
}

func (s *BookingService) AcceptBooking(bookingID uint) error {
	booking, err := s.repo.GetBookingByID(bookingID)
	if err != nil {
		return errors.New("booking not found")
	}

	if booking.Status != "Pending" {
		return errors.New("booking cannot be accepted because it is not in 'Pending' status")
	}

	booking.Status = "Accepted"
	return s.repo.UpdateBooking(booking)
}

func (s *BookingService) DeclineBooking(bookingID uint) error {
	booking, err := s.repo.GetBookingByID(bookingID)
	if err != nil {
		return errors.New("booking not found")
	}

	if booking.Status != "Pending" {
		return errors.New("booking cannot be declined because it is not in 'Pending' status")
	}

	booking.Status = "Declined"
	return s.repo.UpdateBooking(booking)
}


func (s *BookingService) GetBookingsByUserID(userID uint) ([]model.Booking, error) {
	return s.repo.GetBookingsByUserID(userID)
}

func (s *BookingService) GetBookingByID(bookingID uint) (*model.Booking, error) {
	return s.repo.GetBookingByID(bookingID)
}

func (s *BookingService) UpdateBooking(booking *model.Booking) error {
	return s.repo.UpdateBooking(booking)
}

func (s *BookingService) DeleteBooking(bookingID uint) error {
	return s.repo.DeleteBooking(bookingID)
}
