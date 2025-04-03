package repository

import (
	"context"

	"github.com/LAshinCHE/ticket_booking_service/internal/models"
	"github.com/LAshinCHE/ticket_booking_service/internal/repository/schemas"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"

	sq "github.com/Masterminds/squirrel"
)

const (
	bookingTable = "booking"
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

func (r *Repository) GetBookingByID(ctx context.Context, bookingID models.BookingID) (*models.Booking, error) {
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
func ToDomainBooking(booking schemas.Booking) *models.Booking {
	return &models.Booking{
		ID:      models.BookingID(booking.ID),
		UserID:  models.GetUserID(booking.UserID),
		Tikcets: models.TicketID(booking.TicketID),
		Status:  models.MapBookingStatus(booking.Status),
	}
}
