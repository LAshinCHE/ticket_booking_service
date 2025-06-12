package http

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	_ "github.com/LAshinCHE/ticket_booking_service/booking-service/docs"
	"github.com/LAshinCHE/ticket_booking_service/booking-service/internal/api/http/types"
	"github.com/LAshinCHE/ticket_booking_service/booking-service/internal/metrics"
	"github.com/LAshinCHE/ticket_booking_service/booking-service/internal/models"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

//go:generate mockgen -package http -source=http.go -destination http_mocks.go
type (
	service interface {
		CheckTicketIsBooking(ctx context.Context, ticketID int) (bool, error)
		GetBookingByID(ctx context.Context, id int) (*models.Booking, error)
		CreateBooking(ctx context.Context, req models.CreateBookingData) (int, error)
		CreateBookingInternal(ctx context.Context, req models.BookingRequset) (int, error)
		DeleteBookingInternal(ctx context.Context, id int) error
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
	r.HandleFunc("/booking/", handler.CreateBookingHandler).Methods("POST")

	internal := r.PathPrefix("/internal").Subrouter()
	internal.Use(types.InternalAuthMiddleware)
	internal.HandleFunc("/booking/create", handler.CreateBookingInternalHandler).Methods("POST")
	internal.HandleFunc("/booking/delete", handler.DeleteBookingInternalHandler).Methods("DELETE")

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
	_, span := otel.Tracer("booking-service").Start(request.Context(), "HealthCheck")
	defer span.End()
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
	start := time.Now()

	ctx, span := otel.Tracer("booking-service").Start(request.Context(), "GetBookingByIDHandler")
	defer span.End()

	handlerName := "GetBookingByIDHandler"

	defer func() {
		duration := float64(time.Since(start).Milliseconds())
		metrics.ObserveRequestDuration(ctx, handlerName, duration)
	}()

	id, err := types.GetBookingByID(request)
	if err != nil {
		metrics.IncClientErrorByHandler(ctx, handlerName, http.StatusBadRequest)
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	booking, err := h.service.GetBookingByID(ctx, id)
	if err != nil {
		metrics.IncServerErrorByHandler(ctx, handlerName, http.StatusInternalServerError)
		http.Error(writer, internalServerError.Error(), http.StatusInternalServerError)
		return
	}

	metrics.IncOkRespByHandler(ctx, handlerName)
	types.ProcessError(writer, err, &types.GetBookingByIDHandlerResponse{Booking: booking})
}

// CreateBookingHandler godoc
// @Summary Создать бронь
// @Description Создание бронирования пользователем
// @Tags booking
// @Accept json
// @Produce json
// @Param booking body models.CreateBookingData true "Данные для создания брони"
// @Success 200 {object} types.CreateBookingResponse
// @Failure 400 {string} string "bad request"
// @Failure 500 {string} string "internal server error"
// @Router /booking/ [post]
func (h *Handler) CreateBookingHandler(writer http.ResponseWriter, request *http.Request) {
	start := time.Now()
	ctx, span := otel.Tracer("booking-service").Start(request.Context(), "CreateBookingHandler")
	defer span.End()
	handlerName := "CreateBookingHandler"
	metrics.IncRespByHandler(ctx, handlerName)
	defer func() {
		duration := float64(time.Since(start).Milliseconds())
		metrics.ObserveRequestDuration(ctx, handlerName, duration)
	}()

	data, err := types.CreateBooking(request)
	if err != nil {
		metrics.IncClientErrorByHandler(ctx, handlerName, http.StatusBadRequest)
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := h.service.CreateBooking(ctx, data)
	if err != nil {
		metrics.IncServerErrorByHandler(ctx, handlerName, http.StatusInternalServerError)
		http.Error(writer, internalServerError.Error(), http.StatusInternalServerError)
		return
	}

	metrics.IncOkRespByHandler(ctx, handlerName)
	types.ProcessError(writer, err, &types.CreateBookingResponse{BookingID: id})
}

// CreateBookingInternalHandler godoc
// @Summary Внутреннее создание бронирования
// @Description Создание бронирования внутренним сервисом
// @Tags internal
// @Accept json
// @Produce json
// @Param booking body types.CreateBookingInternalRequest true "Данные для создания брони"
// @Success 200 {object} types.CreateBookingResponse
// @Failure 400 {string} string "invalid request"
// @Failure 500 {string} string "internal server error"
// @Router /internal/booking/create [post]
// @Security ApiKeyAuth
func (h *Handler) CreateBookingInternalHandler(w http.ResponseWriter, r *http.Request) {
	var req types.CreateBookingInternalRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	propagator := propagation.TraceContext{}
	ctx := propagator.Extract(r.Context(), propagation.HeaderCarrier(r.Header))

	tracer := otel.Tracer("booking-service")
	ctx, span := tracer.Start(ctx, "CreateBookingInternalHandler")
	defer span.End()

	booking := models.BookingRequset{
		UserID:   req.UserID,
		TicketID: req.TicketID,
	}

	id, err := h.service.CreateBookingInternal(ctx, booking)
	log.Println("[BookingCreateInternal]Booking id is :", id)
	log.Println("[BookingCreateInternal]Service error: ", err)
	types.ProcessError(w, err, &types.CreateBookingResponse{BookingID: id})
}

// DeleteBookingInternalHandler godoc
// @Summary Внутреннее удаление бронирования
// @Description Удаление бронирования внутренним сервисом
// @Tags internal
// @Accept json
// @Produce plain
// @Param booking body types.DeleteBookingInternalRequest true "Данные для удаления брони"
// @Success 200 {string} string "ok"
// @Failure 400 {string} string "invalid request"
// @Failure 500 {string} string "failed to create booking"
// @Router /internal/booking/delete [delete]
// @Security ApiKeyAuth
func (h *Handler) DeleteBookingInternalHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	defer log.Printf("BOOKING SERVICE: DELETE_BOOKING err %s \n", err)
	propagator := propagation.TraceContext{}
	ctx := propagator.Extract(r.Context(), propagation.HeaderCarrier(r.Header))

	tracer := otel.Tracer("booking-service")
	ctx, span := tracer.Start(ctx, "DeleteBookingInternalHandler")
	defer span.End()

	var req types.DeleteBookingInternalRequest
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	bookingID := req.BookingID

	if err = h.service.DeleteBookingInternal(ctx, bookingID); err != nil {
		http.Error(w, "failed to create booking", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

type Handler struct {
	service service
}

var internalServerError = errors.New("internal server error")
