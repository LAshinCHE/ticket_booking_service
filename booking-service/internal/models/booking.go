package models

import "github.com/google/uuid"

type (
	BookingID uuid.UUID
	UserID    uuid.UUID
	TicketID  uuid.UUID
)

type BookingStatus string

const (
	BookingStatusDraft    BookingStatus = "draft"    // Черновик брони
	BookingStatusReserved BookingStatus = "reserved" // Забронировано, но не оплачено
	BookingStatusPaid     BookingStatus = "paid"     // Оплачено
	BookingStatusCanceled BookingStatus = "canceled" // Отменено
)

type Booking struct {
	ID      uuid.UUID
	UserID  uuid.UUID
	Tikcets uuid.UUID
	Status  BookingStatus
}

func MapBookingStatus(status string) BookingStatus {
	switch status {
	case "draft":
		return BookingStatusDraft
	case "reserved":
		return BookingStatusReserved
	case "paid":
		return BookingStatusPaid
	case "canceled":
		return BookingStatusCanceled
	default:
		return BookingStatusDraft
	}
}
