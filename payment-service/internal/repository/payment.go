package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewReopository(db *pgxpool.Pool) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) GetBalance(ctx context.Context, paymentID uuid.UUID) (float64, error) {
	return 0, nil
}

func (r *Repository) UpdateBalance(ctx context.Context, paymentID uuid.UUID, payment float64) error {
	return nil
}
