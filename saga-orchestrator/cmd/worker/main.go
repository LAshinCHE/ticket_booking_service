package main

import (
	"log"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func main() {
	c, err := client.Dial(client.Options{HostPort: "temporal:7233"})
	if err != nil {
		log.Fatalln(err)
	}

	worker := worker.New(c, "booking-saga-task-queue", worker.Options{
		MetricsScope: tally.NewScope("saga", nil), // метрики в Prometheus
	})
	worker.RegisterWorkflow(workflows.BookingSaga)
	worker.RegisterActivity(activities.CheckAvailability /* и др. */)

	log.Fatal(worker.Run(worker.InterruptCh()))
}
