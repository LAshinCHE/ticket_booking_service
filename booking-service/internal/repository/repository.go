package repository

import (
	"context"
	"fmt"
	"log"

	"github.com/LAshinCHE/ticket_booking_service/booking-service/internal/models"
	"github.com/LAshinCHE/ticket_booking_service/booking-service/internal/repository/schemas"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.opentelemetry.io/otel"

	sq "github.com/Masterminds/squirrel"
)

const (
	bookingTable = "bookings"
)

var (
	bookingColumns = []string{"id", "user_id", "ticket_id", "status", "created_at", "updated_at"}
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(driver *pgxpool.Pool) *Repository {
	return &Repository{
		driver,
	}
}

func (r *Repository) GetBookingByID(ctx context.Context, bookingID int) (*models.Booking, error) {
	query := sq.Select(bookingColumns...).
		From(bookingTable).
		Where("id = $1", bookingID).PlaceholderFormat(sq.Dollar)

	rawQuery, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	var booking schemas.Booking

	if err := pgxscan.Select(ctx, r.db, &booking, rawQuery, args...); err != nil {
		return nil, err
	}

	return ToDomainBooking(booking), nil
}

func (r *Repository) CheckTicketIsBooked(ctx context.Context, ticketID int) (bool, error) {
	var exists bool

	query := `
        SELECT EXISTS (
            SELECT 1 FROM bookings 
            WHERE ticket_id = $1 AND status = 'booked'
        )`

	err := r.db.QueryRow(ctx, query, ticketID).Scan(&exists)
	if err != nil {
		return false, err
	}

	return !exists, nil
}

func (r *Repository) CreateBooking(ctx context.Context, booking models.Booking) error {
	ctx, span := otel.Tracer("booking-service").Start(ctx, "Repository.Booking.CreateBooking")
	defer span.End()
	query := sq.Insert(bookingTable).
		Columns("id", "user_id", "ticket_id", "status").
		Values(booking.ID, booking.UserID, booking.TicketID, booking.Status).
		PlaceholderFormat(sq.Dollar)

	queryString, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build insert booking sql: %w", err)
	}

	_, err = r.db.Exec(ctx, queryString, args...)
	if err != nil {
		return fmt.Errorf("failed to execute insert booking: %w", err)
	}

	return nil
}

func (r *Repository) DeleteBookingByID(ctx context.Context, bookingID int) error {
	ctx, span := otel.Tracer("booking-service").Start(ctx, "Repository.Booking.DeleteBookingInternal")
	defer span.End()
	log.Printf("DELETE BOOKING REPOSITORY BOOKING ID: %d \n", bookingID)
	query := sq.Delete(bookingTable).
		Where(sq.Eq{"id": bookingID}).
		PlaceholderFormat(sq.Dollar)

	queryString, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build delete booking sql: %w", err)
	}

	result, err := r.db.Exec(ctx, queryString, args...)
	if err != nil {
		return fmt.Errorf("failed to execute insert booking: %w", err)
	}

	rowAffected := result.RowsAffected()
	if rowAffected == 0 {
		log.Printf("no booking found with id: %d \n", bookingID)
		return fmt.Errorf("no booking found with id: %d", bookingID)
	}

	return nil
}

func (r *Repository) BookingChangeStatus(ctx context.Context, bookingID int, status models.BookingStatus) error {
	ctx, span := otel.Tracer("booking-service").Start(ctx, "Repository.Booking.BookingChangeStatus")
	defer span.End()

	query := sq.Update(bookingTable).
		Set("status", status).
		Where(sq.Eq{"id": bookingID}).
		PlaceholderFormat(sq.Dollar)

	queryString, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("Faild to generate sql query: %w", err)
	}

	result, err := r.db.Exec(ctx, queryString, args...)
	if err != nil {
		return fmt.Errorf("failed to execute booking change status: %w", err)
	}

	rowAffected := result.RowsAffected()
	if rowAffected == 0 {
		log.Printf("no booking found with id: %d \n", bookingID)
		return fmt.Errorf("no booking found with id: %d", bookingID)
	}

	return nil
}

func ToDomainBooking(booking schemas.Booking) *models.Booking {

	return &models.Booking{
		ID:       booking.ID,
		UserID:   booking.UserID,
		TicketID: booking.TicketID,
		Status:   models.MapBookingStatus(booking.Status),
	}
}
