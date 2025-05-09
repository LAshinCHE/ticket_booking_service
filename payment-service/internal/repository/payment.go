package repository

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.opentelemetry.io/otel"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewReopository(db *pgxpool.Pool) *Repository {
	return &Repository{
		db: db,
	}
}

var ErrInsufficientFunds = errors.New("insufficient funds")

func (r *Repository) DebitUserBalance(ctx context.Context, userID int64, amount float64) error {
	ctx, span := otel.Tracer("payment-service").Start(ctx, "Repository.Payment.DebitUserBalance")
	defer span.End()
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("start tx: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			tx.Commit(ctx)
		}
	}()

	var currentBalance float64
	err = tx.QueryRow(ctx, `
		SELECT balance
		FROM users
		WHERE id = $1
		FOR UPDATE
	`, userID).Scan(&currentBalance)
	if err != nil {
		return fmt.Errorf("select balance: %w", err)
	}
	log.Printf("CurentBalance: %f, ammount %f", currentBalance, amount)
	if currentBalance-amount < 0 {
		return fmt.Errorf("Curent balance less than ammount")
	}

	_, err = tx.Exec(ctx, `
		UPDATE users
		SET balance = balance - $1
		WHERE id = $2
	`, amount, userID)
	if err != nil {
		return fmt.Errorf("update balance: %w", err)
	}

	return nil
}

func (r *Repository) CreditUserBalance(ctx context.Context, userID int64, amount float64) error {
	ctx, span := otel.Tracer("payment-service").Start(ctx, "Repository.Payment.CreditUserBalance")
	defer span.End()
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("start tx: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		} else {
			_ = tx.Commit(ctx)
		}
	}()

	_, err = tx.Exec(ctx, `
		UPDATE users
		SET balance = balance + $1
		WHERE id = $2
	`, amount, userID)
	if err != nil {
		return fmt.Errorf("update balance: %w", err)
	}

	return nil
}
