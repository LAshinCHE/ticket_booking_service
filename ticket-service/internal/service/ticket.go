package service

import (
	"context"
	"crypto/rand"
	"errors"
	"math/big"

	"github.com/LAshinCHE/ticket_booking_service/ticket-service/internal/models"
	"go.opentelemetry.io/otel"
)

type RepositoryTicket interface {
	MakeaAvailable(ctx context.Context, id int) error
	GetAvailability(ctx context.Context, ticketID int) (bool, error)
	GetTicket(ctx context.Context, ticketID int) (*models.Ticket, error)
	CreateTicket(ctx context.Context, ticket models.Ticket) error
	ReserveTickert(ctx context.Context, ticketID int) error
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

func (t *Ticket) GetTicket(ctx context.Context, ticketID int) (*models.Ticket, error) {
	ctx, span := otel.Tracer("ticket-service").Start(ctx, "service.Ticket.GetTicket")
	defer span.End()

	return t.RepositoryTicket.GetTicket(ctx, ticketID)
}

func (t *Ticket) ReserveTickert(ctx context.Context, ticketID int) error {
	ctx, span := otel.Tracer("ticket-service").Start(ctx, "service.Ticket.ReserveTickert")
	defer span.End()
	availability, err := t.RepositoryTicket.GetAvailability(ctx, ticketID)
	if err != nil {
		return err
	}
	if !availability {
		return errors.New("Ticket is not available, chose another ticket")
	}

	err = t.RepositoryTicket.ReserveTickert(ctx, ticketID)
	if err != nil {
		return err
	}
	return nil
}

func (t *Ticket) CheckTicket(ctx context.Context, ticketID int) (bool, error) {
	ctx, span := otel.Tracer("ticket-service").Start(ctx, "service.Ticket.CreateTicket")
	defer span.End()
	availability, err := t.RepositoryTicket.GetAvailability(ctx, ticketID)
	if err != nil {
		return false, err
	}

	return availability, nil
}

func (t *Ticket) CreateTicket(ctx context.Context, ticketParam models.TicketModelParamRequest) (int, error) {
	ctx, span := otel.Tracer("ticket-service").Start(ctx, "service.Ticket.CreateTicket")
	defer span.End()
	max := big.NewInt(1000000)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return 0, err
	}
	id := int(n.Int64())
	ticket := models.Ticket{
		ID:        id,
		Price:     ticketParam.Price,
		Available: true,
	}

	return id, t.RepositoryTicket.CreateTicket(ctx, ticket)

}

func (t *Ticket) MackeAvailable(ctx context.Context, ticketID int) error {
	ctx, span := otel.Tracer("ticket-service").Start(ctx, "service.Ticket.MackeAvailable")
	defer span.End()
	return t.RepositoryTicket.MakeaAvailable(ctx, ticketID)
}
