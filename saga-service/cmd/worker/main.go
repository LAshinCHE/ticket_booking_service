package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/LAshinCHE/ticket_booking_service/saga-service/activities"
	sagawf "github.com/LAshinCHE/ticket_booking_service/saga-service/workflow"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func main() {
	log.Println("Try listen service 0.0.0.0:7233")
	tc, err := client.Dial(client.Options{
		HostPort: "temporal:7233",
	})
	if err != nil {
		log.Fatalf("init Temporal client: %v", err)
	}
	defer tc.Close()
	svcs := activities.NewServiceClients(
		"http://localhost:8081",
		"http://localhost:8082",
		"http://localhost:8083",
		"http://localhost:8084",
	)
	acts := activities.NewBookingActivities(svcs)

	w := worker.New(tc, "BOOKING_SAGA_QUEUE", worker.Options{})

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
