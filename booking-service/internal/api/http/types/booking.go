package types

import (
	"encoding/json"
	"net/http"

	"github.com/LAshinCHE/ticket_booking_service/booking-service/internal/models"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

const SagaServiceToken = "super-secure-saga-token"

type GetBookingByIDHandlerResponse struct {
	Booking *models.Booking
}

type CreateBookingResponse struct {
	BookingID uuid.UUID
}

type CreateBookingInternalRequest struct {
	BookingID string `json:"booking_id"`
	UserID    int64  `json:"user_id"`
	TicketID  string `json:"ticket_id"`
}

type DeleteBookingInternalRequest struct {
	BookingID string `json:"booking_id"`
}

func GetBookingByID(r *http.Request) (uuid.UUID, error) {
	vars := mux.Vars(r)
	uuidStr, ok := vars["booking_id"]
	if !ok || len(uuidStr) == 0 {
		return uuid.Nil, MissingUUID
	}

	id, err := uuid.Parse(uuidStr)
	if err != nil {
		return uuid.Nil, InvalidUUID
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
