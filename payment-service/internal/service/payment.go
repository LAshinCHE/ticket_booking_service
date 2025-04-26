package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/LAshinCHE/ticket_booking_service/payment-service/internal/repository"
)

type PaymentRepository interface {
	DebitUserBalance(ctx context.Context, userID int64, amount float64) error
	CreditUserBalance(ctx context.Context, userID int64, amount float64) error
}

type Deps struct {
	PaymentRepository
}

type Service struct {
	Deps
}

func NewPaymentService(d Deps) *Service {
	return &Service{
		d,
	}
}

func (s *Service) DebitBalance(ctx context.Context, userID int64, amount float64) (bool, error) {
	err := s.PaymentRepository.DebitUserBalance(ctx, userID, amount)
	if err != nil {
		if errors.Is(err, repository.ErrInsufficientFunds) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (s *Service) RefundBalance(ctx context.Context, userID int64, amount float64) error {
	if amount <= 0 {
		return fmt.Errorf("refund amount must be positive")
	}
	err := s.PaymentRepository.CreditUserBalance(ctx, userID, amount)
	if err != nil {
		return err
	}
	return nil
}

var ErrInsufficientFunds = errors.New("insufficient funds")
