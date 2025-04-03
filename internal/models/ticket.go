package models

type (
	TicketID int64
)

type TicketStatus string

const (
	TicketActiveStatus   TicketStatus = "active"
	TicketInactiveStatus TicketStatus = "inactive"
	TicketArchivedStatus TicketStatus = "archived"
	TicketDeletedStatus  TicketStatus = "deleted"
)

type Ticket struct {
	id     int64
	status TicketStatus
}
