package schemas

import (
	"database/sql"
	"time"
)

type Booking struct { // структура бронирования приложения
	ID        int64         `db:"id"`
	UserID    sql.NullInt64 `db:"user_id"`
	TicketID  int64         `db:"ticket_id"`
	Status    string        `db:"status"`
	CreatedAt time.Time     `db:"created_at"`
	UpdatedAt time.Time     `db:"updated_at"`
}
