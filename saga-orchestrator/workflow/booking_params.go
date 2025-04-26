package workflow

type BookingParams struct {
	UserID   int64   `json:"user_id"`
	TicketID string  `json:"ticket_id"`
	Price    float64 `json:"price"`
}
