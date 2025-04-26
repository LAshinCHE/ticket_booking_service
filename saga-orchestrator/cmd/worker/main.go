package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/LAshinCHE/ticket_booking_service/saga-orchestrator/activities"
	"github.com/LAshinCHE/ticket_booking_service/saga-orchestrator/clients"
	sagawf "github.com/LAshinCHE/ticket_booking_service/saga-orchestrator/workflow"
	"go.temporal.io/sdk/worker"
)

func main() {
	tc, err := clients.NewTemporalClient()
	if err != nil {
		log.Fatalf("init Temporal client: %v", err)
	}
	defer tc.Client.Close()

	svcs := activities.NewServiceClients(
		"http://booking:8080",
		"http://ticket:8081",
		"http://payment:8082",
		"http://notify:8083",
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
