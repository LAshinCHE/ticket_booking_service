package schemas

import (
	"time"

	"github.com/google/uuid"
)

type Booking struct { // структура бронирования приложения
	ID        string        `db:"id"`
	UserID    uuid.NullUUID `db:"user_id"`
	TicketID  string        `db:"ticket_id"`
	Status    string        `db:"status"`
	CreatedAt time.Time     `db:"created_at"`
	UpdatedAt time.Time     `db:"updated_at"`
}
