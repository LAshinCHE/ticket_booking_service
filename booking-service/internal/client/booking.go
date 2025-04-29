package client

import (
	"context"
	"time"

	"github.com/LAshinCHE/ticket_booking_service/booking-service/internal/models"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.temporal.io/sdk/client"
)

type BookingWorkflowInput struct {
	BookingData models.CreateBookingData
	TraceCtx    map[string]string
}

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
	ctx, span := otel.Tracer("booking-service").Start(ctx, "StartBookingSaga")
	defer span.End()
	options := client.StartWorkflowOptions{
		ID:        "booking_workflow_" + time.Now().Format("20060102150405"),
		TaskQueue: "BOOKING_SAGA_TASK_QUEUE",
	}

	propagator := propagation.TraceContext{}
	traceMap := make(map[string]string)
	propagator.Inject(ctx, propagation.MapCarrier(traceMap))

	workflowInput := BookingWorkflowInput{
		BookingData: req,
		TraceCtx:    traceMap,
	}

	_, err := s.Client.ExecuteWorkflow(ctx, options, "BookingWorkflow", workflowInput)
	return err

}
