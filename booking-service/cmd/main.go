package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"
	"time"

	"github.com/LAshinCHE/ticket_booking_service/booking-service/cmd/internal"
	internalhttp "github.com/LAshinCHE/ticket_booking_service/booking-service/internal/api/http"
	"github.com/LAshinCHE/ticket_booking_service/booking-service/internal/client"
	"github.com/LAshinCHE/ticket_booking_service/booking-service/internal/domain/service"
	"github.com/LAshinCHE/ticket_booking_service/booking-service/internal/metrics"
	"github.com/LAshinCHE/ticket_booking_service/booking-service/internal/repository"
	"github.com/LAshinCHE/ticket_booking_service/booking-service/internal/tracer"
)

const (
	applicationPort = ":8081"

	watchDuration    = 3 * time.Second
	shutdownDuration = 5 * time.Second
)

// @title Booking Service API
// @version 1.0
// @description API для сервиса бронирования

// @host localhost:8081
// @BasePath /

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	db, err := internal.NewPostgreSQL()
	if err != nil {
		log.Fatalf("Could not initialize Database connection %s", err)
	}
	defer db.Close()

	tracer.MustSetup(ctx, "booking-service")
	metrics.InitMetrics()
	//metrics.StartMetricsEndpoint()

	repoBooking := repository.NewRepository(db)

	temporalClient, err := client.NewTemporalClient()
	if err != nil {
		log.Fatalf("Could not initialize temporal client to interapt with saga service %s", err)
	}
	defer temporalClient.Client.Close()

	bookingService := service.NewBookingService(service.Deps{
		RepositoryBooking: repoBooking,
		SagaClient:        temporalClient,
	})

	internalhttp.MustRun(ctx, shutdownDuration, applicationPort, bookingService)

}
