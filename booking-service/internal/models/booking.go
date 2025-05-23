package models

import (
	"github.com/google/uuid"
)

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
	ID       int
	UserID   int
	TicketID int
	Status   BookingStatus
}

type BookingRequset struct {
	UserID   int           `json:"userID"`
	TicketID int           `json:"ticketID"`
	Status   BookingStatus `json:"status"`
}

type CreateBookingData struct {
	Price    float64 `json:"price"`
	TicketID int     `json:"ticketID"`
	UserID   int     `json:"userID"`
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
