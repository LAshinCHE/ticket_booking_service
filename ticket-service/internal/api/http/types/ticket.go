package types

import (
	"encoding/json"
	"net/http"

	"github.com/LAshinCHE/ticket_booking_service/ticket-service/internal/models"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func GetTicketIDRequest(r *http.Request) (uuid.UUID, error) {
	vars := mux.Vars(r)
	uuidStr, ok := vars["ticket_id"]
	if !ok || len(uuidStr) == 0 {
		return uuid.Nil, BadUUID
	}
	id, err := uuid.Parse(uuidStr)

	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

func CreateTicketRequest(r *http.Request) (models.TicketModelParamRequest, error) {
	var req models.TicketModelParamRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&req)
	if err != nil {
		return models.TicketModelParamRequest{}, err
	}
	return req, nil
}
