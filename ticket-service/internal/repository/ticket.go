package repository

import (
	"context"
	"fmt"

	"github.com/LAshinCHE/ticket_booking_service/ticket-service/internal/models"
	"github.com/LAshinCHE/ticket_booking_service/ticket-service/internal/repository/schemas"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"

	sq "github.com/Masterminds/squirrel"
)

const (
	ticketTable = "tickets"
)

var (
	ticketColumns = []string{"id", "price", "available"}
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(driver *pgxpool.Pool) *Repository {
	return &Repository{
		driver,
	}
}

func (r *Repository) GetAvailability(ctx context.Context, ticketID int) (bool, error) {
	ctx, span := otel.Tracer("ticket-service").Start(ctx, "Repository.Ticket.GetAvailability")
	defer span.End()
	query := sq.Select(ticketColumns...).
		From(ticketTable).
		Where("id = $1", ticketID).PlaceholderFormat(sq.Dollar)

	rawQuery, args, err := query.ToSql()
	if err != nil {
		return false, err
	}

	var ticket schemas.Ticket

	if err := pgxscan.Select(ctx, r.db, &ticket, rawQuery, args...); err != nil {
		return false, err
	}

	return ticket.Available, nil
}

func (r *Repository) GetTicket(ctx context.Context, ticketID int) (*models.Ticket, error) {

	ctx, span := otel.Tracer("ticket-service").Start(ctx, "Repository.Ticket.GetTicket")
	defer span.End()
	query := sq.Select(ticketColumns...).
		From(ticketTable).
		Where("id = $1", ticketID).
		PlaceholderFormat(sq.Dollar)

	rawQuery, args, err := query.ToSql()
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to build SQL query")
		return nil, err
	}

	var ticket schemas.Ticket
	if err := pgxscan.Get(ctx, r.db, &ticket, rawQuery, args...); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "ticket not found or db error")
		return nil, fmt.Errorf("ticket with id %v not found, with error: %v ", args, err)
	}

	span.SetStatus(codes.Ok, "ticket successfully retrieved")
	return ToDomainTicket(ticket), nil
}

func (r *Repository) MakeaAvailable(ctx context.Context, ticketID int) error {
	ctx, span := otel.Tracer("ticket-service").Start(ctx, "Repository.Ticket.MakeaAvailable")
	defer span.End()
	query := sq.Update(ticketTable).
		Set("available", true).
		Where("id = $1", ticketID).
		PlaceholderFormat(sq.Dollar)

	rawQuery, args, err := query.ToSql()
	if err != nil {
		return err
	}

	result, err := r.db.Exec(ctx, rawQuery, args...)
	if err != nil {
		return fmt.Errorf("failed to execute query: %w", err)
	}

	rowsAffected := result.RowsAffected()

	if rowsAffected == 0 {
		return fmt.Errorf("ticket with id %d not found", ticketID)
	}

	return nil
}

func (r *Repository) CreateTicket(ctx context.Context, ticket models.Ticket) error {
	ctx, span := otel.Tracer("ticket-service").Start(ctx, "Repository.Ticket.CreateTicket")
	defer span.End()
	query := sq.Insert(ticketTable).
		Columns(ticketColumns...).
		Values(ticket.ID, ticket.Price, ticket.Available).
		PlaceholderFormat(sq.Dollar)

	rawQuery, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}
	_, err = r.db.Exec(ctx, rawQuery, args...)
	if err != nil {
		return fmt.Errorf("failed to execute insert: %w", err)
	}

	return nil
}

func (r *Repository) ReserveTickert(ctx context.Context, ticketID int) error {
	ctx, span := otel.Tracer("ticket-service").Start(ctx, "Repository.Ticket.ReserveTickert")
	defer span.End()
	query := sq.Update(ticketTable).Set("available", false).Where(sq.Eq{"id": ticketID}).PlaceholderFormat(sq.Dollar)

	rawQuery, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}
	_, err = r.db.Exec(ctx, rawQuery, args...)
	if err != nil {
		return fmt.Errorf("failed to execute update: %w", err)
	}
	return nil
}

func ToDomainTicket(ticket schemas.Ticket) *models.Ticket {

	return &models.Ticket{
		ID:        ticket.ID,
		Price:     ticket.Price,
		Available: ticket.Available,
	}
}
