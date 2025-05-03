package schemas

type Ticket struct {
	ID        int     `db:"id"`
	Price     float64 `db:"price"`
	Available bool    `db:"available"`
}
