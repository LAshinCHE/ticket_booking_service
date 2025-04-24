package models

type BookingParams struct {
	ID       int64
	UserID   int64
	TicketID int64
	Amount   float64
}
