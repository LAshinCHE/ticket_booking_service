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

	//imdb - базовая реализация по кешам включающая в себя два поля мапку и лог
	// telegramClient := client.NewTelegram()
	// slackClient := client.NewSlack()
	// whatsappClient := client.NewWhatsapp()
	// contactsClient := client.NewContacts(ctx)

	// contactService := contact.NewNotificationPosition(contactsClient, positionIMDB)

	// notifiers := notification.New(telegramClient, slackClient, whatsappClient)
	// leClient := client.NewLeaderElectioneer()

	// positionWatcher := watcher.NewPositionWatcher(watchDuration, positionIMDB, notifiers, leClient)
	// go positionWatcher.Run(ctx)

	// internalhttp.MustRun(ctx, shutdownDuration, applicationPort, positionIMDB)

}
