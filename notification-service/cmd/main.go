package main

import (
	"context"
	"os/signal"
	"syscall"
	"time"

	myhttp "github.com/LAshinCHE/ticket_booking_service/notification-service/internal/api/http"
	"github.com/LAshinCHE/ticket_booking_service/notification-service/internal/service"
	"github.com/LAshinCHE/ticket_booking_service/notification-service/internal/tracer"
)

var (
	shutdownDuration = 5 * time.Second
	applicationPort  = ":8084"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	tracer.MustSetup(ctx, "notification-service")
	service := service.NewNotificationService()

	myhttp.MustRun(ctx, applicationPort, service, shutdownDuration)
}
