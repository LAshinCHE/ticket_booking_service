package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/LAshinCHE/ticket_booking_service/saga-service/activities"
	"github.com/LAshinCHE/ticket_booking_service/saga-service/metrics"
	tracer "github.com/LAshinCHE/ticket_booking_service/saga-service/trace"
	sagawf "github.com/LAshinCHE/ticket_booking_service/saga-service/workflow"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tracer.MustSetup(ctx, "saga-worker")

	tc, err := client.Dial(client.Options{
		HostPort: "temporal:7233",
	})
	if err != nil {
		log.Fatalf("init Temporal client: %v", err)
	}
	defer tc.Close()

	svcs := activities.NewServiceClients(
		"http://booking-service:8081",
		"http://ticket-service:8082",
		"http://payment-service:8083",
		"http://notification-service:8084",
	)

	if err := metrics.Init(ctx, "otel-collector:4317", "saga-service"); err != nil {
		log.Fatalf("metrics init: %v", err)
	}

	defer metrics.Shutdown(ctx)

	acts := activities.NewBookingActivities(svcs)
	w := worker.New(tc, "BOOKING_SAGA_QUEUE", worker.Options{})
	log.Println("Registrate worker")

	w.RegisterWorkflowWithOptions(sagawf.BookingSagaWorkflow, workflow.RegisterOptions{
		Name: "BookingWorkflow",
	})
	w.RegisterActivity(acts)
	log.Println("Registrate worckflow")

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	stop := make(chan interface{})

	go func() {
		<-sig
		close(stop)
		cancel()
	}()

	if err := w.Run(stop); err != nil {
		log.Fatal("worker stopped:", err)
	}
}
