package main

import (
	"log"

	"github.com/LAshinCHE/ticket_booking_service/saga-orchestrator/activities"
	"github.com/LAshinCHE/ticket_booking_service/saga-orchestrator/clients"
	workflows "github.com/LAshinCHE/ticket_booking_service/saga-orchestrator/workflow"
	"go.temporal.io/sdk/worker"
)

func main() {
	temporalClient, err := clients.NewTemporalClient()
	if err != nil {
		log.Fatalf("Could not initialize temporal client to interapt with saga service %s", err)
	}
	defer temporalClient.Client.Close()

	w := worker.New(temporalClient.Client, "BOOKING_SAGA_QUEUE", worker.Options{})

	w.RegisterWorkflow(workflows.BookingSagaWorkflow)
	w.RegisterActivity(activities.WorkflowActivitie{})

	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("unable to start Worker", err)
	}

}
