package schemas

import "github.com/google/uuid"

type Ticket struct {
	ID        uuid.UUID `db:"id"`
	Price     float64   `db:"price"`
	Available bool      `db:"available"`
}
