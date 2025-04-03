package service

import (
	"context"

	"github.com/LAshinCHE/ticket_booking_service/internal/models"
)

type Repository interface {
	GetBookingByID(ctx context.Context, id models.BookingID) (*models.Booking, error)
}

type Deps struct {
	Repository
}

type Booking struct {
	Deps
}

func NewBookingService(d Deps) *Booking {
	return &Booking{
		d,
	}
}

func (b *Booking) GetBookingByID(ctx context.Context, id models.BookingID) (*models.Booking, error) {
	return b.Repository.GetBookingByID(ctx, id)
}

func (b *Booking) CreateBooking(ctx context.Context, userID models.UserID, ticket models.UserID) error {
	return nil
}
