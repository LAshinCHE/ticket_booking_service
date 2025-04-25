package main

import (
	"log"

	workflows "github.com/LAshinCHE/ticket_booking_service/saga-orchestrator/workflow"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func main() {
	c, err := client.Dial(client.Options{HostPort: "temporal:7233"})
	if err != nil {
		log.Fatalln(err)
	}

	worker := worker.New(c, "booking-saga-task-queue", worker.Options{
		MetricsScope: tally.NewScope("saga", nil),
	})
	worker.RegisterWorkflow(workflows.BookingSagaWorkflow())
	worker.RegisterActivity(activities.CheckAvailability /* и др. */)

	log.Fatal(worker.Run(worker.InterruptCh()))
}
