package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/LAshinCHE/ticket_booking_service/payment-service/cmd/internal"
	myhttp "github.com/LAshinCHE/ticket_booking_service/payment-service/internal/api/http"
	"github.com/LAshinCHE/ticket_booking_service/payment-service/internal/repository"
	"github.com/LAshinCHE/ticket_booking_service/payment-service/internal/service"
)

var (
	shutdownDuration = 5 * time.Second
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)

	defer cancel()

	db, err := internal.NewPostgreSQL()
	if err != nil {
		log.Fatal(err)
	}
	repository := repository.NewReopository(db)

	service := service.NewPaymentService(service.Deps{
		PaymentRepository: repository,
	})

	myhttp.MustRun(ctx, service, os.Getenv("PORT"), shutdownDuration)
}
