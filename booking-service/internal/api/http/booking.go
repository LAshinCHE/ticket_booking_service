package http

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	_ "github.com/LAshinCHE/ticket_booking_service/booking-service/docs"
	"github.com/LAshinCHE/ticket_booking_service/booking-service/internal/api/http/types"
	"github.com/LAshinCHE/ticket_booking_service/booking-service/internal/models"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

//go:generate mockgen -package http -source=http.go -destination http_mocks.go
type (
	service interface {
		CheckTicketIsBooking(ctx context.Context, ticketID uuid.UUID) (bool, error)
		GetBookingByID(ctx context.Context, id uuid.UUID) (*models.Booking, error)
	}
)

func MustRun(ctx context.Context, shutdownDur time.Duration, addr string, app service) {
	handler := &Handler{
		service: app,
	}

	// mux := http.NewServeMux()
	r := mux.NewRouter()
	r.HandleFunc("/", handler.HealthCheck).Methods("GET")
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)
	r.HandleFunc("/booking/{booking_id}", handler.GetBookingByIDHandler).Methods("GET")

	server := &http.Server{
		Addr:    addr,
		Handler: r,
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

// HealthCheck godoc
// @Summary Health check
// @Description Проверка доступности сервиса
// @Tags health
// @Produce plain
// @Success 200 {string} string "hello"
// @Router / [get]
func (h *Handler) HealthCheck(writer http.ResponseWriter, request *http.Request) {
	writer.Write([]byte("hello"))
}

// GetBookingByIDHandler godoc
// @Summary Получить бронь по ID
// @Description Возвращает информацию о бронировании
// @Tags booking
// @Produce json
// @Param booking_id path string true "Booking ID"
// @Success 200 {object} types.GetBookingByIDHandlerResponse
// @Failure 400 {string} string "bad request"
// @Failure 500 {string} string "internal server error"
// @Router /booking/{booking_id} [get]
func (h *Handler) GetBookingByIDHandler(writer http.ResponseWriter, request *http.Request) {
	uuid, err := types.GetBookingByID(request)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
	}

	booking, err := h.service.GetBookingByID(request.Context(), uuid)

	types.ProcessError(writer, err, &types.GetBookingByIDHandlerResponse{Booking: booking})
}

// func (h *Handler) CreateBookingHandler(writer http.ResponseWriter, request *http.Request) {
// 	uuid, err := types.CreateBooking(request)
// 	if err != nil {
// 		http.Error(writer, err.Error(), http.StatusBadRequest)
// 	}

// 	booking, err := h.service.GetBookingByID(request.Context(), uuid)

// 	types.ProcessError(writer, err, &types.GetBookingByIDHandlerResponse{Booking: booking})

// }

type Handler struct {
	service service
}

var internalServerError = errors.New("internal server error")
