package http

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/LAshinCHE/ticket_booking_service/ticket-service/docs"
	"github.com/LAshinCHE/ticket_booking_service/ticket-service/internal/api/http/types"
	"github.com/LAshinCHE/ticket_booking_service/ticket-service/internal/models"
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

	r := mux.NewRouter()
	r.HandleFunc("/", handler.HealthCheck).Methods("GET")
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)
	r.HandleFunc("/ticket/{ticket_id}", handler.GetTicketByIDHandler).Methods("GET")
	r.HandleFunc("/ticket/check/{ticket_id}", handler.CheckTicketHandler).Methods("GET")
	r.HandleFunc("/ticket/", handler.CreateTicketHandler).Methods("POST")

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
	ticketID, err := types.GetTicketIDRequest(r)
	fmt.Printf("Ticket handler Get %v by request", ticketID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	ticket, err := h.service.GetTicket(r.Context(), ticketID)

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
	params, err := types.CreateTicketRequest(r)
	if err != nil || params.Price == 0 {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	id, err := h.service.CreateTicket(r.Context(), params)

	types.ProcessError(w, err, &types.CreateTicketResponse{Id: id})
}

type Handler struct {
	service service
}
