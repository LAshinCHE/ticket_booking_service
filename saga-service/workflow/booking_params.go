package workflow

type CreateBookingData struct {
	Price    int `json:"price"`
	TicketID int `json:"ticketID"`
	UserID   int `json:"userID"`
}

type BookingWorkflowInput struct {
	BookingData CreateBookingData
	TraceCtx    map[string]string
}
