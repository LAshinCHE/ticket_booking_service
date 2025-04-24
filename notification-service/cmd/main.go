package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	myhttp "github.com/LAshinCHE/ticket_booking_service/notification-service/internal/api/http"
	"github.com/LAshinCHE/ticket_booking_service/notification-service/internal/service"
)

var (
	shutdownDuration = 5 * time.Second
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	service := service.NewNotificationService()
	addr := os.Getenv("SERVICE_ADDRES")

	myhttp.MustRun(ctx, addr, service, shutdownDuration)
}
