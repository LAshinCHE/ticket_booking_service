package schemas

type Ticket struct {
	ID      int64   `db:"id"`
	EventID int64   `db:"event_id"`
	Price   float64 `db:"price"`
	Seat    int64   `db:"seat"`
	Status  string  `db:"status"`
}
