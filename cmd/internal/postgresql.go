package internal

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPostgreSQL() (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		return nil, err
	}
	return pool, nil
}
