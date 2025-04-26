package service

import (
	"context"
	"errors"

	"github.com/LAshinCHE/ticket_booking_service/ticket-service/internal/models"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

type RepositoryTicket interface {
	MakeaAvailable(ctx context.Context, id uuid.UUID) error
	GetAvailability(ctx context.Context, ticketID uuid.UUID) (bool, error)
	GetTicket(ctx context.Context, ticketID uuid.UUID) (*models.Ticket, error)
	CreateTicket(ctx context.Context, ticket models.Ticket) error
	UpdateTicketAvaible(ctx context.Context, ticketID uuid.UUID) error
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

func (t *Ticket) GetTicket(ctx context.Context, ticketID uuid.UUID) (*models.Ticket, error) {
	ctx, span := otel.Tracer("ticket-service").Start(ctx, "service.GetTicket")
	defer span.End()

	span.SetAttributes(attribute.String("ticket.id", ticketID.String()))

	return t.RepositoryTicket.GetTicket(ctx, ticketID)
}

func (t *Ticket) ReserveTicket(ctx context.Context, ticketID uuid.UUID) error {

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

func (t *Ticket) CheckTicket(ctx context.Context, ticketID uuid.UUID) (bool, error) {
	availability, err := t.RepositoryTicket.GetAvailability(ctx, ticketID)
	if err != nil {
		return false, err
	}

	return availability, nil
}

func (t *Ticket) CreateTicket(ctx context.Context, ticketParam models.TicketModelParamRequest) (uuid.UUID, error) {
	id := uuid.New()
	ticket := models.Ticket{
		ID:        id,
		Price:     ticketParam.Price,
		Available: true,
	}

	return id, t.RepositoryTicket.CreateTicket(ctx, ticket)

}

func (t *Ticket) UpdateTicketAvaible(ctx context.Context, ticketID uuid.UUID) error {
	return t.RepositoryTicket.UpdateTicketAvaible(ctx, ticketID)
}
