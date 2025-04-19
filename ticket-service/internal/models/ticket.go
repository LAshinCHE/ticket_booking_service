package models

import "github.com/google/uuid"

type TicketStatus string

const (
	TicketActiveStatus   TicketStatus = "active"   // доступный для бронирования
	TicketInactiveStatus TicketStatus = "inactive" // не доступный для бронирования
	TicketArchivedStatus TicketStatus = "archived" // находится в архиве
	TicketDeletedStatus  TicketStatus = "deleted"  // удаленн
)

type Ticket struct {
	ID        uuid.UUID
	Price     float64
	Available bool
}
