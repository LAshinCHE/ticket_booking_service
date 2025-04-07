package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/LAshinCHE/ticket_booking_service/internal/models"
)

type RepositoryBooking interface {
	GetBookingByID(ctx context.Context, id models.BookingID) (*models.Booking, error)
	CheckTicketIsBooked(ctx context.Context, ticketID models.TicketID) (bool, error)
}

type RepositoryTicket interface {
	MakeaAvailable(ctx context.Context, id models.TicketID) error
	GetAvailability(ctx context.Context, ticketID models.TicketID) (bool, error)
}

type Deps struct {
	RepositoryBooking
	RepositoryTicket
}

type Booking struct {
	Deps
}

func NewBookingService(d Deps) *Booking {
	return &Booking{
		d,
	}
}

func (b *Booking) CheckTicketIsBooking(ctx context.Context, ticketID models.TicketID) (bool, error) {
	canBook, err := b.RepositoryBooking.CheckTicketIsBooked(ctx, ticketID)
	if err != nil {
		return false, err
	}
	if !canBook {
		return false, fmt.Errorf("Билет свободен, нельзя бронировать")
	}
	return true, nil
}

func (b *Booking) GetBookingByID(ctx context.Context, id models.BookingID) (*models.Booking, error) {
	return b.RepositoryBooking.GetBookingByID(ctx, id)
}

func (b *Booking) MakeBooking(ctx context.Context, ticketID models.TicketID, user_id models.UserID, price float64) error {
	//основная ручка которая выполняет запрос сразу в несколько сервисов
	// проверка забронирован ли у нас билет
	// err := b.CheckTicket() // логика что билет у нас может быть забронирван кем-то другим
	//1. Проверяет доступность билета и резервирует его
	err := b.ReserveTicket(ctx, ticketID) // билет не доступен по каким-то другим причинам условно удален админом
	if err != nil {
		return err
	}
	//2. Проверка суммы у пользака
	// bool, error := b.CheckPaiment
	//3. списание средств err := b.DebitingMoney
	//4.  уведомление err := b.notify
	return nil
}

// Будет методом ticket service
func (b *Booking) ReserveTicket(ctx context.Context, ticketID models.TicketID) error {

	availability, err := b.RepositoryTicket.GetAvailability(ctx, ticketID)
	if err != nil {
		return err
	}
	if !availability {
		return errors.New("Ticket is not available, chose another ticket")
	}

	err = b.RepositoryTicket.MakeaAvailable(ctx, ticketID)
	if err != nil {
		return err
	}
	return nil
}
