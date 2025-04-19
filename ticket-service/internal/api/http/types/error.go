package types

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/LAshinCHE/ticket_booking_service/ticket-service/internal/models"
	"github.com/LAshinCHE/ticket_booking_service/ticket-service/internal/repository"
)

var (
	BadUUID = errors.New("No such uuid")
)

type GetTicketByIDHandlerResponse struct {
	Ticket *models.Ticket
}

type CheckTicketHandlerResponse struct {
	Valid bool
}

func ProcessError(w http.ResponseWriter, err error, resp any) {
	if err == repository.NotFound {
		http.Error(w, "Key not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	if resp != nil {
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, "Internal Error", http.StatusInternalServerError)
		}
	}
}
