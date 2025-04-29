package main

import (
	"context"
	"log"

	"github.com/LAshinCHE/ticket_booking_service/saga-service/workflow"
	"github.com/google/uuid"
	"go.temporal.io/sdk/client"
)

func main() {
	workflowInput := workflow.BookingWorkflowInput{
		Params: workflow.BookingParams{
			ID:       uuid.New(),
			UserID:   1,
			TicketID: uuid.New(),
			Price:    12.0,
		},
		TraceCtx: make(map[string]string),
	}

	c, err := client.Dial(client.Options{
		HostPort: "temporal:7233",
	})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	options := client.StartWorkflowOptions{
		ID:        "testing-estimate-age-example",
		TaskQueue: "BOOKING_SAGA_QUEUE",
	}

	we, err := c.ExecuteWorkflow(context.Background(), options, workflow.BookingSagaWorkflow, workflowInput)
	if err != nil {
		log.Fatalln("Unable to execute workflow", err)
	}

	var result string
	err = we.Get(context.Background(), &result)
	if err != nil {
		log.Fatalln("Unable get workflow result", err)
	}

	log.Println("Workflow result:", result)
}
