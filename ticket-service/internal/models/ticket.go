package models

type (
	TicketID int64
)

type TicketStatus string

const (
	TicketActiveStatus   TicketStatus = "active"   // доступный для бронирования
	TicketInactiveStatus TicketStatus = "inactive" // не доступный для бронирования
	TicketArchivedStatus TicketStatus = "archived" // находится в архиве
	TicketDeletedStatus  TicketStatus = "deleted"  // удаленн
)

type Ticket struct {
	ID        int64
	Price     float64
	Available bool
}
