package models

type UserID int64
type TicketID int64

type BookingSagaParams struct {
	ID       UserID
	TicketID TicketID
	Amount   float64
}
