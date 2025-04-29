package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/LAshinCHE/ticket_booking_service/saga-service/activities"
	"github.com/LAshinCHE/ticket_booking_service/saga-service/clients"
	sagawf "github.com/LAshinCHE/ticket_booking_service/saga-service/workflow"
	"go.temporal.io/sdk/worker"
)

func main() {
	tc, err := clients.NewTemporalClient()
	if err != nil {
		log.Fatalf("init Temporal client: %v", err)
	}
	defer tc.Client.Close()

	svcs := activities.NewServiceClients(
		"http://localhost:8080",
		"http://localhost:8082",
		"http://localhost:8083",
		"http://localhost:8084",
	)
	acts := activities.NewBookingActivities(svcs)

	w := worker.New(tc.Client, "BOOKING_SAGA_QUEUE", worker.Options{})

	w.RegisterWorkflow(sagawf.BookingSagaWorkflow)
	w.RegisterActivity(acts)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	stop := make(chan interface{})

	go func() {
		<-sig
		close(stop)
	}()

	if err := w.Run(stop); err != nil {
		log.Fatal("worker stopped:", err)
	}
}
