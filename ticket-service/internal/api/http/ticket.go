package http

import (
	"context"
	"log"
	"net/http"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"

	_ "github.com/LAshinCHE/ticket_booking_service/ticket-service/docs"
	"github.com/LAshinCHE/ticket_booking_service/ticket-service/internal/api/http/types"
	"github.com/LAshinCHE/ticket_booking_service/ticket-service/internal/models"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

type service interface {
	GetTicket(ctx context.Context, ticketID int) (*models.Ticket, error)
	CheckTicket(ctx context.Context, ticketID int) (bool, error)
	CreateTicket(ctx context.Context, ticketParam models.TicketModelParamRequest) (int, error)

	MackeAvailable(ctx context.Context, ticketID int) error
	ReserveTickert(ctx context.Context, ticketID int) error
}

func MustRun(ctx context.Context, shutdownDur time.Duration, addr string, app service) {
	handler := Handler{
		service: app,
	}

	r := mux.NewRouter()

	r.Use(otelmux.Middleware("ticket-service"))

	r.HandleFunc("/", handler.HealthCheck).Methods("GET")
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)
	r.HandleFunc("/ticket/", handler.CreateTicketHandler).Methods("POST")
	r.HandleFunc("/ticket/{ticket_id}", handler.GetTicketByIDHandler).Methods("GET")
	r.HandleFunc("/ticket/{ticket_id}/check", handler.CheckTicketHandler).Methods("GET")
	// r.HandleFunc("/tickets/{ticket_id}", handler.UpdateTicketAvaibleHandler).Methods("PUT")

	r.HandleFunc("/ticket/{ticket_id}/reserve", handler.ReservTicketHandler).Methods("PUT")
	r.HandleFunc("/ticket/{ticket_id}/available", handler.MakeAvailableHandler).Methods("PUT")

	server := http.Server{
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
// @Summary Проверка здоровья сервиса
// @Description Возвращает простое сообщение
// @Tags health
// @Success 200 {string} string "Hello"
// @Router / [get]
func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello"))
}

func (h *Handler) MakeAvailableHandler(w http.ResponseWriter, r *http.Request) {
	propagator := propagation.TraceContext{}
	ctx := propagator.Extract(r.Context(), propagation.HeaderCarrier(r.Header))

	tracer := otel.Tracer("tikcet-service")
	ctx, span := tracer.Start(ctx, "MakeAvailableHandler")
	defer span.End()

	ticketid, err := types.GetTicketIDRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	err = h.service.MackeAvailable(ctx, ticketid)
	types.ProcessError(w, err, nil)
}

func (h *Handler) ReservTicketHandler(w http.ResponseWriter, r *http.Request) {
	propagator := propagation.TraceContext{}
	ctx := propagator.Extract(r.Context(), propagation.HeaderCarrier(r.Header))

	tracer := otel.Tracer("tikcet-service")
	ctx, span := tracer.Start(ctx, "ReservTicketHandler")
	defer span.End()
	log.Println("ReservTicketHandler")
	ticketid, err := types.GetTicketIDRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	err = h.service.ReserveTickert(ctx, ticketid)
	types.ProcessError(w, err, nil)
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
	ctx, span := otel.Tracer("ticket-service").Start(r.Context(), "GetTicketByIDHandler")
	defer span.End()

	ticketID, err := types.GetTicketIDRequest(r)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "invalid ticket ID")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ticket, err := h.service.GetTicket(ctx, ticketID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to get ticket")
	}
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
	ctx, span := otel.Tracer("ticket-service").Start(r.Context(), "GetTicketByIDHandler")
	defer span.End()

	params, err := types.CreateTicketRequest(r)
	if err != nil || params.Price == 0 {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ticketID, err := h.service.CreateTicket(ctx, params)
	types.ProcessError(w, err, &types.CreateTicketResponse{Id: ticketID})
}

// func (h *Handler) UpdateTicketAvaibleHandler(w http.ResponseWriter, r *http.Request) {
// 	params, id, err := types.UpdateTicketAvaibleRequest(r)
// 	if err != nil || params.Available == false {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}

// 	err = h.service.UpdateTicketAvaible(r.Context(), id)
// 	types.ProcessError(w, err, nil)
// }

type Handler struct {
	service service
}
