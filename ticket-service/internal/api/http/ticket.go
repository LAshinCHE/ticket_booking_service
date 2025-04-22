package http

import (
	"context"
	"log"
	"net/http"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"

	_ "github.com/LAshinCHE/ticket_booking_service/ticket-service/docs"
	"github.com/LAshinCHE/ticket_booking_service/ticket-service/internal/api/http/types"
	"github.com/LAshinCHE/ticket_booking_service/ticket-service/internal/models"
	"github.com/LAshinCHE/ticket_booking_service/ticket-service/internal/tracing"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

type service interface {
	ReserveTicket(ctx context.Context, ticketID uuid.UUID) error
	GetTicket(ctx context.Context, ticketID uuid.UUID) (*models.Ticket, error)
	CheckTicket(ctx context.Context, ticketID uuid.UUID) (bool, error)
	CreateTicket(ctx context.Context, ticketParam models.TicketModelParamRequest) (uuid.UUID, error)
}

func MustRun(ctx context.Context, shutdownDur time.Duration, addr string, app service) {
	handler := Handler{
		service: app,
	}

	tp, err := tracing.InitTracer(ctx, "ticket-service")
	if err != nil {
		log.Fatalf("failed to initialize tracer: %v", err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/", handler.HealthCheck).Methods("GET")
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)
	r.HandleFunc("/ticket/{ticket_id}", handler.GetTicketByIDHandler).Methods("GET")
	r.HandleFunc("/ticket/check/{ticket_id}", handler.CheckTicketHandler).Methods("GET")
	r.HandleFunc("/ticket/", handler.CreateTicketHandler).Methods("POST")

	otelRouter := otelhttp.NewHandler(r, "http-server")
	server := http.Server{
		Addr:    addr,
		Handler: otelRouter,
	}

	go func() {
		<-ctx.Done()

		log.Printf("Shuting down server with duration %0.3fs", shutdownDur.Seconds())
		<-time.After(shutdownDur)

		if err := tp.Shutdown(ctx); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}

		if err := server.Shutdown(context.Background()); err != nil {
			log.Fatalf("Http handler Shutdown: %s", err)
		}
	}()

	if err := server.ListenAndServe(); err != nil {
		log.Printf("HTTP server ListenAndServe: %s", err)
	}

}

// HealthCheck godoc
// @Summary Проверка здоровья сервиса
// @Description Возвращает простое сообщение
// @Tags health
// @Success 200 {string} string "Hello"
// @Router / [get]
func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello"))
}

// CheckTicketHandler godoc
// @Summary Проверка валидности билета
// @Description Проверяет доступность билета
// @Tags ticket
// @Param ticket_id path string true "Ticket UUID"
// @Produce json
// @Success 200 {object} types.CheckTicketHandlerResponse
// @Failure 400 {string} string "bad request"
// @Failure 500 {string} string "internal error"
// @Router /ticket/check/{ticket_id} [get]
func (h *Handler) CheckTicketHandler(w http.ResponseWriter, r *http.Request) {
	ticketID, err := types.GetTicketIDRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	valid, err := h.service.CheckTicket(r.Context(), ticketID)

	types.ProcessError(w, err, &types.CheckTicketHandlerResponse{Valid: valid})
}

// GetTicketByIDHandler godoc
// @Summary Получение информации о билете
// @Description Получает билет по UUID
// @Tags ticket
// @Param ticket_id path string true "Ticket UUID"
// @Produce json
// @Success 200 {object} types.GetTicketByIDHandlerResponse
// @Failure 400 {string} string "bad request"
// @Failure 500 {string} string "internal error"
// @Router /ticket/{ticket_id} [get]
func (h *Handler) GetTicketByIDHandler(w http.ResponseWriter, r *http.Request) {
	tr := otel.Tracer("ticket-service")
	ctx, span := tr.Start(r.Context(), "GetTicketByIDHandler")
	defer span.End()

	ticketID, err := types.GetTicketIDRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ticket, err := h.service.GetTicket(ctx, ticketID)
	types.ProcessError(w, err, &types.GetTicketByIDHandlerResponse{Ticket: ticket})
}

// CreateTicketHandler godoc
// @Summary Создает новый билет
// @Description Создает новый билет по переданным параметрам
// @Tags ticket
// @Param ticket  body      models.TicketModelParamRequest  true  "Данные билета"
// @Accept       json
// @Produce      json
// @Success  200 {object}  types.CreateTicketResponse
// @Failure 400 {string} string "bad request"
// @Failure 500 {string} string "internal error"
// @Router /ticket/ [post]
func (h *Handler) CreateTicketHandler(w http.ResponseWriter, r *http.Request) {
	tr := otel.Tracer("ticket-service")
	ctx, span := tr.Start(r.Context(), "CreateTicketHandler")
	defer span.End()

	params, err := types.CreateTicketRequest(r)
	if err != nil || params.Price == 0 {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := h.service.CreateTicket(ctx, params)
	types.ProcessError(w, err, &types.CreateTicketResponse{Id: id})
}

type Handler struct {
	service service
}
