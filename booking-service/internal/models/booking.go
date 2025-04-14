package models

type (
	BookingID int64
	UserID    int64
	TicketID  int64
)

type BookingStatus string

const (
	BookingStatusDraft    BookingStatus = "draft"    // Черновик брони
	BookingStatusReserved BookingStatus = "reserved" // Забронировано, но не оплачено
	BookingStatusPaid     BookingStatus = "paid"     // Оплачено
	BookingStatusCanceled BookingStatus = "canceled" // Отменено
)

type Booking struct {
	ID      int64
	UserID  int64
	Tikcets int64
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
