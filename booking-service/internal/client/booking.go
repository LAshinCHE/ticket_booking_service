package client

import (
	"context"
	"time"

	"github.com/LAshinCHE/ticket_booking_service/booking-service/internal/models"
	"go.temporal.io/sdk/client"
)

type SagaClient struct {
	Client client.Client
}

func NewTemporalClient() (*SagaClient, error) {
	TemporalClient, err := client.Dial(client.Options{
		HostPort: "temporal:7233",
	})
	if err != nil {
		return nil, err
	}
	return &SagaClient{Client: TemporalClient}, nil
}

func (s *SagaClient) StartBookingSaga(ctx context.Context, req models.CreateBookingData) error {
	options := client.StartWorkflowOptions{
		ID:        "booking_workflow_" + time.Now().Format("20060102150405"),
		TaskQueue: "BOOKING_SAGA_TASK_QUEUE",
	}

	_, err := s.Client.ExecuteWorkflow(ctx, options, "BookingWorkflow", req)
	return err

}
