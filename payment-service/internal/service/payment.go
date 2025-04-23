package service

import (
	"context"

	"github.com/google/uuid"
)

type PaymentRepository interface {
	GetBalance(ctx context.Context, paymentID uuid.UUID) (float64, error)
	UpdateBalance(ctx context.Context, paymentID uuid.UUID, funds float64) error
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

func (s *Service) DebitFundsFromBalance(ctx context.Context, paymentID uuid.UUID, funds float64) error {
	if funds > 0 {
		funds *= -1
	}

	return s.PaymentRepository.UpdateBalance(ctx, paymentID, funds)
}

func (s *Service) RefundFundsToBalance(ctx context.Context, paymentID uuid.UUID, funds float64) error {
	if funds < 0 {
		funds *= -1
	}

	return s.PaymentRepository.UpdateBalance(ctx, paymentID, funds)
}

func (s *Service) GetBalance(ctx context.Context, paymentID uuid.UUID) (float64, error) {
	return s.PaymentRepository.GetBalance(ctx, paymentID)
}
