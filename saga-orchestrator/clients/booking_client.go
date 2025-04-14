package clients

import (
	"net/http"

	"github.com/LAshinCHE/ticket_booking_service/saga-orchestrator/models"
)

type BookingClient struct {
	HTTPClient *http.Client
}

func NewBookingClient() *BookingClient {
	return &BookingClient{
		HTTPClient: &http.Client{},
	}
}

func Reserve(id models.TicketID) (bool, error) {
	return false, nil
}
