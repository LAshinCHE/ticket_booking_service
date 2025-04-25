package activities

import (
	"context"
	"fmt"
	"os"

	"github.com/LAshinCHE/ticket_booking_service/saga-orchestrator/clients"
	"github.com/LAshinCHE/ticket_booking_service/saga-orchestrator/models"
)

type WorkflowActivitie struct {
}

func (w *WorkflowActivitie) CheckBookingActivity(ctx context.Context, params models.BookingParams) error {
	bookingClient := clients.NewBookingClient(os.Getenv("BOOKING_ADDRES"))
	err := bookingClient.CheckBooking(ctx, params.ID)
	if err != nil {
		return fmt.Errorf("CheckBookingActivity failed: %w", err)
	}

	return nil
}
