package schemas

type Ticket struct {
	ID        int64   `db:"id"`
	Price     float64 `db:"price"`
	Available bool    `db:"available"`
}
