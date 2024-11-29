package handler

import (
	"encoding/json"
	"net/http"
	"os"
	"rentora-go/internal/middleware"
	"rentora-go/internal/model"
	"rentora-go/internal/service"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type BookingHandler struct {
	service *service.BookingService
}

func NewBookingHandler(service *service.BookingService) *BookingHandler {
	return &BookingHandler{service: service}
}

func (h *BookingHandler) CreateBooking(w http.ResponseWriter, r *http.Request) {
	var booking model.Booking
	if err := json.NewDecoder(r.Body).Decode(&booking); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Here you might add logic to check car availability, user validation, etc.
	if err := h.service.CreateBooking(&booking); err != nil {
		http.Error(w, "Error creating booking", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(booking)
}

func (h *BookingHandler) GetBookingsByUserID(w http.ResponseWriter, r *http.Request) {
	userIDStr := chi.URLParam(r, "userID")
	userID, err := strconv.ParseUint(userIDStr, 10, 32) // Parse string to uint64
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	bookings, err := h.service.GetBookingsByUserID(uint(userID)) // Convert to uint
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(bookings)
}

func (h *BookingHandler) GetBookingByID(w http.ResponseWriter, r *http.Request) {
	bookingIDStr := chi.URLParam(r, "bookingID")
	bookingID, err := strconv.ParseUint(bookingIDStr, 10, 32) // Parse string to uint64
	if err != nil {
		http.Error(w, "Invalid booking ID", http.StatusBadRequest) // Fix error message
		return
	}

	// Cast to uint before passing to service
	booking, err := h.service.GetBookingByID(uint(bookingID))
	if err != nil {
		http.Error(w, "Booking not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(booking)
}

func (h *BookingHandler) UpdateBooking(w http.ResponseWriter, r *http.Request) {
	var booking model.Booking
	if err := json.NewDecoder(r.Body).Decode(&booking); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if err := h.service.UpdateBooking(&booking); err != nil {
		http.Error(w, "Error updating booking", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(booking)
}

func (h *BookingHandler) DeleteBooking(w http.ResponseWriter, r *http.Request) {
	bookingIDStr := chi.URLParam(r, "bookingID")
	bookingID, err := strconv.ParseUint(bookingIDStr, 10, 32) // Parse string to uint64
	if err != nil {
		http.Error(w, "Invalid booking ID", http.StatusBadRequest) // Fix error message
		return
	}

	if err := h.service.DeleteBooking(uint(bookingID)); err != nil {
		http.Error(w, "Error deleting booking", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}


func (h *BookingHandler) AcceptBooking(w http.ResponseWriter, r *http.Request) {
	bookingIDStr := chi.URLParam(r, "bookingID")
	bookingID, err := strconv.ParseUint(bookingIDStr, 10, 32) // Parse string to uint64
	if err != nil {
		http.Error(w, "Invalid booking ID", http.StatusBadRequest)
		return
	}

	err = h.service.AcceptBooking(uint(bookingID)) // Convert to uint
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Booking accepted"})
}



func (h *BookingHandler) DeclineBooking(w http.ResponseWriter, r *http.Request) {
	bookingIDStr := chi.URLParam(r, "bookingID")
	bookingID, err := strconv.ParseUint(bookingIDStr, 10, 32) // Parse string to uint64
	if err != nil {
		http.Error(w, "Invalid booking ID", http.StatusBadRequest)
		return
	}

	err = h.service.DeclineBooking(uint(bookingID)) // Convert to uint
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Booking declined"})
}



func RegisterBookingRoutes(r chi.Router, bookingHandler *BookingHandler) {
	jwtSecret := []byte(os.Getenv("JWT_SECRET"))
    r.Group(func(public chi.Router) {
        public.Post("/bookings", bookingHandler.CreateBooking)      // Create booking
        public.Get("/bookings/{bookingID}", bookingHandler.GetBookingByID) // Get booking by ID
    })

    r.Group(func(protected chi.Router) {
        protected.Use(middleware.AuthMiddleware(jwtSecret))                   // Apply authentication middleware
        protected.Get("/bookings", bookingHandler.GetBookingsByUserID) // Get bookings for the user
        protected.Put("/bookings/{bookingID}/accept", bookingHandler.AcceptBooking) // Accept booking
        protected.Put("/bookings/{bookingID}/decline", bookingHandler.DeclineBooking) // Decline booking
        protected.Delete("/bookings/{bookingID}", bookingHandler.DeleteBooking) // Delete booking
    })

    r.Get("/bookings/health", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("Bookings service is running"))
    })
}