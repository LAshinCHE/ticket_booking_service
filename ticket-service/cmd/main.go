package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"
	"time"

	"github.com/LAshinCHE/ticket_booking_service/ticket-service/cmd/internal"
	"github.com/LAshinCHE/ticket_booking_service/ticket-service/internal/repository"
	"github.com/LAshinCHE/ticket_booking_service/ticket-service/internal/service"

	internalhttp "github.com/LAshinCHE/ticket_booking_service/ticket-service/internal/api/http"
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
		log.Fatal(err)
	}

	defer db.Close()

	repositoryTicket := repository.NewRepository(db)

	app := service.NewBookingService(service.Deps{
		RepositoryTicket: repositoryTicket,
	})

	internalhttp.MustRun(ctx, shutdownDuration, applicationPort, app)
}
