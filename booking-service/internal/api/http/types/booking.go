package types

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/LAshinCHE/ticket_booking_service/booking-service/internal/models"
	"github.com/gorilla/mux"
)

const SagaServiceToken = "super-secure-saga-token"

type GetBookingByIDHandlerResponse struct {
	Booking *models.Booking
}

type CreateBookingResponse struct {
	BookingID int `json:"bookingID"`
}

type CreateBookingInternalRequest struct {
	UserID   int `json:"user_id"`
	TicketID int `json:"ticket_id"`
}

type DeleteBookingInternalRequest struct {
	BookingID int `json:"booking_id"`
}

func GetBookingByID(r *http.Request) (int, error) {
	vars := mux.Vars(r)
	idStr, ok := vars["booking_id"]
	if !ok {
		return 0, errors.New("[Bad request]: Missing id")
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func InternalAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token != SagaServiceToken {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func CreateBooking(r *http.Request) (models.CreateBookingData, error) {
	var req models.CreateBookingData
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return models.CreateBookingData{}, err
	}

	return req, nil
}
