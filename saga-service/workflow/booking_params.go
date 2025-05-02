package workflow

import "github.com/google/uuid"

type CreateBookingData struct {
	ID       uuid.UUID `json:"id"`
	Price    int       `json:"price"`
	TicketID uuid.UUID `json:"ticketID"`
	UserID   int       `json:"userID"`
}

type BookingWorkflowInput struct {
	BookingData CreateBookingData
	TraceCtx    map[string]string
}
