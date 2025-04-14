package service

import (
	"context"
	"errors"

	"github.com/LAshinCHE/ticket_booking_service/ticket-service/internal/models"
)

type RepositoryTicket interface {
	MakeaAvailable(ctx context.Context, id models.TicketID) error
	GetAvailability(ctx context.Context, ticketID models.TicketID) (bool, error)
}

type Deps struct {
	RepositoryTicket
}

type Ticket struct {
	Deps
}

func NewBookingService(d Deps) *Ticket {
	return &Ticket{
		d,
	}
}

func (t *Ticket) GetTicket(ctx context.Context, ticketID models.TicketID) {

}

func (t *Ticket) ReserveTicket(ctx context.Context, ticketID models.TicketID) error {

	availability, err := t.RepositoryTicket.GetAvailability(ctx, ticketID)
	if err != nil {
		return err
	}
	if !availability {
		return errors.New("Ticket is not available, chose another ticket")
	}

	err = t.RepositoryTicket.MakeaAvailable(ctx, ticketID)
	if err != nil {
		return err
	}
	return nil
}
