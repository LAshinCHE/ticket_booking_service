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
	"github.com/LAshinCHE/ticket_booking_service/ticket-service/internal/tracer"

	internalhttp "github.com/LAshinCHE/ticket_booking_service/ticket-service/internal/api/http"
)

const (
	applicationPort = ":8082"

	watchDuration    = 3 * time.Second
	shutdownDuration = 5 * time.Second
)

// @title Ticket Service API
// @version 1.0
// @description API для сервиса бронирования

// @host localhost:8082
// @BasePath /

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	db, err := internal.NewPostgreSQL()
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	tracer.MustSetup(ctx, "ticket-service")
	repositoryTicket := repository.NewRepository(db)

	app := service.NewBookingService(service.Deps{
		RepositoryTicket: repositoryTicket,
	})

	internalhttp.MustRun(ctx, shutdownDuration, applicationPort, app)
}
