package types

import (
	"encoding/json"
	"net/http"

	"github.com/LAshinCHE/ticket_booking_service/booking-service/internal/models"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type GetBookingByIDHandlerResponse struct {
	Booking *models.Booking
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

func CreateBooking(r *http.Request) (models.CreateBookingData, error) {
	var req models.CreateBookingData
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return models.CreateBookingData{}, err
	}

	return req, nil
}
