package workflow

import "github.com/google/uuid"

type BookingParams struct {
	ID       uuid.UUID
	UserID   int64
	TicketID uuid.UUID
	Price    float64
}

type BookingWorkflowInput struct {
	Params   BookingParams
	TraceCtx map[string]string
}
