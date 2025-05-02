package service

import (
	"context"
	"fmt"
	"log"

	"github.com/LAshinCHE/ticket_booking_service/booking-service/internal/models"
	"go.opentelemetry.io/otel"
)

type RepositoryBooking interface {
	GetBookingByID(ctx context.Context, id int) (*models.Booking, error)
	CheckTicketIsBooked(ctx context.Context, ticketID int) (bool, error)
	CreateBooking(ctx context.Context, booking models.Booking) error
	DeleteBookingByID(ctx context.Context, bookingID int) error
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
func (b *Booking) CheckTicketIsBooking(ctx context.Context, ticketID int) (bool, error) {
	canBook, err := b.RepositoryBooking.CheckTicketIsBooked(ctx, ticketID)
	if err != nil {
		return false, err
	}
	if !canBook {
		return false, fmt.Errorf("Билет свободен, нельзя бронировать")
	}
	return true, nil
}

func (b *Booking) GetBookingByID(ctx context.Context, id int) (*models.Booking, error) {
	return b.RepositoryBooking.GetBookingByID(ctx, id)
}

func (b *Booking) CreateBookingInternal(ctx context.Context, booking models.Booking) error {
	ctx, span := otel.Tracer("booking-service").Start(ctx, "Service.Booking.CreateBookingInternal")
	defer span.End()
	booking.Status = models.BookingStatusDraft
	return b.RepositoryBooking.CreateBooking(ctx, booking)
}

func (b *Booking) DeleteBookingInternal(ctx context.Context, bookingID int) error {
	return b.RepositoryBooking.DeleteBookingByID(ctx, bookingID)
}

func (b *Booking) CreateBooking(ctx context.Context, req models.CreateBookingData) (int, error) {
	ctx, span := otel.Tracer("booking-service").Start(ctx, "BookingService.CreateBooking")
	defer span.End()

	if err := b.SagaClient.StartBookingSaga(ctx, req); err != nil {
		log.Printf("failed to start booking saga: %v", err)
		return 0, err
	}

	return req.ID, nil
}
