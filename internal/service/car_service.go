package service

import (
	"rentora-go/internal/model"
	"rentora-go/internal/repository"

)

type CarService struct {
	repo repository.CarRepository
}

func NewCarService(repo repository.CarRepository) *CarService {
	return &CarService{repo: repo}
}

func (s *CarService) CreateCarListing(car *model.Car) error {
	return s.repo.CreateCar(car)
}

func (s *CarService) GetAvailableCars(location string) ([]model.Car, error) {
	return s.repo.GetAvailableCars(location)
}

func (s *CarService) UpdateCarListing(car *model.Car) error {
	return s.repo.UpdateCar(car)
}

func (s *CarService) DeleteCarListing(carID uint) error {
	return s.repo.DeleteCar(carID)
}
