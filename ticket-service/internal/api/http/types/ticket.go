package types

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/LAshinCHE/ticket_booking_service/ticket-service/internal/models"
	"github.com/gorilla/mux"
)

func GetTicketIDRequest(r *http.Request) (int, error) {
	vars := mux.Vars(r)
	idstr, ok := vars["ticket_id"]
	if !ok {
		return 0, fmt.Errorf("No ticket id data in request string")
	}
	id, err := strconv.Atoi(idstr)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func UpdateTicketAvaibleRequest(r *http.Request) (models.TicketUpdateAvaibleData, int, error) {
	vars := mux.Vars(r)
	idStr, ok := vars["ticket_id"]
	if !ok {
		return models.TicketUpdateAvaibleData{}, 0, fmt.Errorf("cant get id from request string")
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return models.TicketUpdateAvaibleData{}, 0, err
	}

	var req models.TicketUpdateAvaibleData
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&req)
	if err != nil {
		return models.TicketUpdateAvaibleData{}, 0, err
	}
	return req, id, nil
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
