package http

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/LAshinCHE/ticket_booking_service/booking-service/internal/models"
)

//go:generate mockgen -package http -source=http.go -destination http_mocks.go
type (
	booking interface {
		CheckTicketIsBooking(ctx context.Context, ticketID models.TicketID) (bool, error)
		GetBookingByID(ctx context.Context, id models.BookingID) (*models.Booking, error)
	}
)

func MustRun(ctx context.Context, shutdownDur time.Duration, addr string, app booking) {
	handler := &Handler{
		booking: app,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", handler.Hello) // Changed to HandleFunc
	mux.HandleFunc("/booking", handler.GetBookingByID)

	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	go func() {
		<-ctx.Done()

		log.Printf("Shuting down server with duration %0.3fs", shutdownDur.Seconds())
		<-time.After(shutdownDur)

		if err := server.Shutdown(context.Background()); err != nil {
			log.Fatalf("Http handler Shutdown: %s", err)
		}
	}()

	if err := server.ListenAndServe(); err != nil {
		log.Printf("HTTP server ListenAndServe: %s", err)
	}
}

func (h *Handler) Hello(writer http.ResponseWriter, request *http.Request) {
	writer.Write([]byte("hello"))
}

func (h *Handler) GetBookingByID(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("GetBookingByID func handler")
	// Implementation here
}

type Handler struct {
	booking booking
}

var internalServerError = errors.New("internal server error")
