package main

import (
	"context"
	"os/signal"
	"syscall"
	"time"

	myhttp "github.com/LAshinCHE/ticket_booking_service/notification-service/internal/api/http"
	"github.com/LAshinCHE/ticket_booking_service/notification-service/internal/service"
)

var (
	shutdownDuration = 5 * time.Second
	applicationPort  = "0.0.0.0:8084"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	service := service.NewNotificationService()

	myhttp.MustRun(ctx, applicationPort, service, shutdownDuration)
}
