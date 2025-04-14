package service

import (
	"context"
	"fmt"

	"github.com/LAshinCHE/ticket_booking_service/booking-service/internal/models"
)

type RepositoryBooking interface {
	GetBookingByID(ctx context.Context, id models.BookingID) (*models.Booking, error)
	CheckTicketIsBooked(ctx context.Context, ticketID models.TicketID) (bool, error)
	CreateBooking(ctx context.Context, booking *models.Booking) (int64, error)
}

type SagaClient interface {
	StartBookingSaga(ctx context.Context, bookingID int64, userID int64, ticketID int64) error
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

// Создает бронь (создает отправляет статус запроса в saga-service) в случае ошибки saga-service принимает решение о дальнейших действиях
func (b *Booking) CreateBooking(ctx context.Context, userID int64, ticketID int64) (int64, error) {
	booking := &models.Booking{
		UserID:  userID,
		Tikcets: ticketID,
		Status:  models.BookingStatusDraft,
	}

	bookingID, err := b.RepositoryBooking.CreateBooking(ctx, booking)
	if err != nil {
		return 0, err
	}

	err = b.SagaClient.StartBookingSaga(ctx, bookingID, userID, ticketID)
	if err != nil {
		return bookingID, fmt.Errorf("failed to start saga: %w", err)
	}

	return bookingID, nil
}

// func (b *Booking) MakeBooking(ctx context.Context, ticketID models.TicketID, user_id models.UserID, price float64) error {
// 	//основная ручка которая выполняет запрос сразу в несколько сервисов
// 	// проверка забронирован ли у нас билет
// 	// err := b.CheckTicket() // логика что билет   нас может быть забронирван кем-то другим
// 	//1. Проверяет доступность билета и резервирует его
// 	err := b.ReserveTicket(ctx, ticketID) // билет не доступен по каким-то другим причинам условно удален админом
// 	if err != nil {
// 		return err
// 	}
// 	//2. Проверка суммы у пользака
// 	// bool, error := b.CheckPaiment
// 	//3. списание средств err := b.DebitingMoney
// 	//4.  уведомление err := b.notify
// 	return nil
// }
