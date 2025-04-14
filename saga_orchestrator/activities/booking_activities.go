package activities

import (
	"context"
	"fmt"

	"github.com/LAshinCHE/ticket_booking_service/saga-orchestrator/models"
)

type BookingClient interface {
	Reserve(id models.TicketID) (bool, error)
}

func ReserveTicketActivity(ctx context.Context, params models.BookingSagaParams) error {
	// Пример вызова ticket-service через HTTP или gRPC
	ok, err := BookingClient.Reserve(params.ID)
	if err != nil || !ok {
		return fmt.Errorf("не удалось зарезервировать билет: %w", err)
	}
	return nil
}
