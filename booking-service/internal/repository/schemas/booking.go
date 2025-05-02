package schemas

import (
	"time"
)

type Booking struct { // структура бронирования приложения
	ID        int       `db:"id"`
	UserID    int       `db:"user_id"`
	TicketID  int       `db:"ticket_id"`
	Status    string    `db:"status"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
