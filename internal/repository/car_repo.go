package repository

import (
	"rentora-go/internal/model"
	"gorm.io/gorm"
)

type CarRepository interface {
	CreateCar(car *model.Car) error
	GetAvailableCars(location string) ([]model.Car, error)
	UpdateCar(car *model.Car) error
	DeleteCar(carID uint) error
}

type carRepository struct {
	db *gorm.DB
}

func NewCarRepository(db *gorm.DB) CarRepository {
	return &carRepository{db: db}
}

func (r *carRepository) CreateCar(car *model.Car) error {
	if err := r.db.Create(car).Error; err != nil {
		return err
	}
	return nil
}

func (r *carRepository) GetAvailableCars(location string) ([]model.Car, error) {
	var cars []model.Car
	if err := r.db.Where("availability = ? AND location = ?", true, location).Find(&cars).Error; err != nil {
		return nil, err
	}
	return cars, nil
}

func (r *carRepository) UpdateCar(car *model.Car) error {
	if err := r.db.Save(car).Error; err != nil {
		return err
	}
	return nil
}

func (r *carRepository) DeleteCar(carID uint) error {
	if err := r.db.Delete(&model.Car{}, carID).Error; err != nil {
		return err
	}
	return nil
}
