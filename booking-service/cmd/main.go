package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"
	"time"

	"github.com/LAshinCHE/ticket_booking_service/booking-service/cmd/internal"
	internalhttp "github.com/LAshinCHE/ticket_booking_service/booking-service/internal/api/http"
	"github.com/LAshinCHE/ticket_booking_service/booking-service/internal/domain/service"
	"github.com/LAshinCHE/ticket_booking_service/booking-service/internal/repository"
)

const (
	applicationPort = ":8080"

	watchDuration    = 3 * time.Second
	shutdownDuration = 5 * time.Second
)

// @title Booking Service API
// @version 1.0
// @description API для сервиса бронирования

// @host localhost:8080
// @BasePath /

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	db, err := internal.NewPostgreSQL()
	if err != nil {
		log.Fatalf("Could not initialize Database connection %s", err)
	}
	defer db.Close()

	repoBooking := repository.NewRepository(db)

	bookingService := service.NewBookingService(service.Deps{
		RepositoryBooking: repoBooking,
	})

	internalhttp.MustRun(ctx, shutdownDuration, applicationPort, bookingService)

}
