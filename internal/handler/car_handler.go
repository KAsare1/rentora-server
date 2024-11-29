package handler

import (
	"encoding/json"
	"net/http"
	"rentora-go/internal/model"
	"rentora-go/internal/service"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type CarHandler struct {
	service *service.CarService
}

func NewCarHandler(service *service.CarService) *CarHandler {
	return &CarHandler{service: service}
}

// RegisterCarRoutes registers the car routes with the router.
func RegisterCarRoutes(r chi.Router, handler *CarHandler) {
	r.Post("/cars", handler.CreateCar)
	r.Get("/cars", handler.GetCars)
	r.Put("/cars/{id}", handler.UpdateCar)
	r.Delete("/cars/{id}", handler.DeleteCar)
}

func (h *CarHandler) CreateCar(w http.ResponseWriter, r *http.Request) {
	var car model.Car
	if err := json.NewDecoder(r.Body).Decode(&car); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if err := h.service.CreateCarListing(&car); err != nil {
		http.Error(w, "Failed to create car listing", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(car)
}

func (h *CarHandler) GetCars(w http.ResponseWriter, r *http.Request) {
	location := r.URL.Query().Get("location")
	cars, err := h.service.GetAvailableCars(location)
	if err != nil {
		http.Error(w, "Failed to retrieve car listings", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(cars)
}

func (h *CarHandler) UpdateCar(w http.ResponseWriter, r *http.Request) {
	carIDStr := chi.URLParam(r, "id")
	carID, err := strconv.ParseUint(carIDStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid car ID", http.StatusBadRequest)
		return
	}

	var car model.Car
	if err := json.NewDecoder(r.Body).Decode(&car); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Set the car ID from the parsed value
	car.ID = uint(carID)

	if err := h.service.UpdateCarListing(&car); err != nil {
		http.Error(w, "Failed to update car listing", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(car)
}

func (h *CarHandler) DeleteCar(w http.ResponseWriter, r *http.Request) {
	carIDStr := chi.URLParam(r, "id")
	carID, err := strconv.ParseUint(carIDStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid car ID", http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteCarListing(uint(carID)); err != nil {
		http.Error(w, "Failed to delete car listing", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
