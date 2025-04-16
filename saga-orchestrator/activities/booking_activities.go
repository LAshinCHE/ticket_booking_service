package activities

import (
	"context"
	"fmt"

	"github.com/LAshinCHE/ticket_booking_service/saga-orchestrator/clients"
	"github.com/LAshinCHE/ticket_booking_service/saga-orchestrator/models"
)

func ReserveTicketActivity(ctx context.Context, params models.BookingSagaParams) error {
	// Пример вызова ticket-service через HTTP или gRPC
	ok, err := clients.BookingClient.Reserve(params.ID)
	if err != nil || !ok {
		return fmt.Errorf("не удалось зарезервировать билет: %w", err)
	}
	return nil
}
