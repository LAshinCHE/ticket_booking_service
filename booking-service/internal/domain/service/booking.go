package service

import (
	"context"
	"fmt"
	"log"

	"github.com/LAshinCHE/ticket_booking_service/booking-service/internal/models"
	"github.com/google/uuid"
)

type RepositoryBooking interface {
	GetBookingByID(ctx context.Context, id uuid.UUID) (*models.Booking, error)
	CheckTicketIsBooked(ctx context.Context, ticketID uuid.UUID) (bool, error)
	CreateBooking(ctx context.Context, booking models.Booking) error
	DeleteBookingByID(ctx context.Context, bookingID uuid.UUID) error
}

type SagaClient interface {
	StartBookingSaga(ctx context.Context, prams models.CreateBookingData) error
}

type Deps struct {
	RepositoryBooking
	SagaClient
}

type Booking struct {
	Deps
}

func NewBookingService(d Deps) *Booking {
	return &Booking{
		d,
	}
}

// Функция которая проверяет не забронирован ли билет
func (b *Booking) CheckTicketIsBooking(ctx context.Context, ticketID uuid.UUID) (bool, error) {
	canBook, err := b.RepositoryBooking.CheckTicketIsBooked(ctx, ticketID)
	if err != nil {
		return false, err
	}
	if !canBook {
		return false, fmt.Errorf("Билет свободен, нельзя бронировать")
	}
	return true, nil
}

func (b *Booking) GetBookingByID(ctx context.Context, id uuid.UUID) (*models.Booking, error) {
	return b.RepositoryBooking.GetBookingByID(ctx, id)
}

func (b *Booking) CreateBookingInternal(ctx context.Context, booking models.Booking) error {
	booking.Status = models.BookingStatusDraft
	return b.RepositoryBooking.CreateBooking(ctx, booking)
}

func (b *Booking) DeleteBookingInternal(ctx context.Context, bookingID uuid.UUID) error {
	return b.RepositoryBooking.DeleteBookingByID(ctx, bookingID)
}

func (b *Booking) CreateBooking(ctx context.Context, req models.CreateBookingData) error {
	req.ID = uuid.New()

	if err := b.SagaClient.StartBookingSaga(ctx, req); err != nil {
		log.Printf("failed to start booking saga: %v", err)
		return err
	}

	return nil
}
