package client

import (
	"context"
	"fmt"
	"log"
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

type TemporalClient struct {
	Client client.Client
}

func NewTemporalClient() (*TemporalClient, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	c, err := client.DialContext(ctx, client.Options{
		HostPort: "temporal:7233",
	})
	if err != nil {
		log.Fatalf("Unable to create Temporal client: %v", err)
	}
	if err != nil {
		return nil, err
	}
	return &TemporalClient{Client: c}, nil
}

func (s *TemporalClient) StartBookingSaga(ctx context.Context, req models.CreateBookingData) (int, error) {
	ctx, span := otel.Tracer("booking-service").Start(ctx, "SagaClient.StartBookingSaga")
	defer span.End()
	options := client.StartWorkflowOptions{
		ID:        "booking_workflow_" + time.Now().Format("20060102150405"),
		TaskQueue: "BOOKING_SAGA_QUEUE",
	}

	propagator := propagation.TraceContext{}
	traceMap := make(map[string]string)
	propagator.Inject(ctx, propagation.MapCarrier(traceMap))

	workflowInput := BookingWorkflowInput{
		BookingData: req,
		TraceCtx:    traceMap,
	}

	fmt.Println("start worckflow", traceMap)
	we, err := s.Client.ExecuteWorkflow(ctx, options, "BookingWorkflow", workflowInput)
	if err != nil {
		log.Fatalln("Unable to execute workflow", err)
	}
	var bookingID int
	err = we.Get(ctx, &bookingID)
	if err != nil {
		log.Println("Workflow execution failed:", err)
		return -1, err
	}

	fmt.Println("Workflow completed successfully")

	return bookingID, nil
}
