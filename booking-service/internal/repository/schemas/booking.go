package schemas

import (
	"time"

	"github.com/google/uuid"
)

type Booking struct { // структура бронирования приложения
	ID        uuid.UUID `db:"id"`
	UserID    uuid.UUID `db:"user_id"`
	TicketID  uuid.UUID `db:"ticket_id"`
	Status    string    `db:"status"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
