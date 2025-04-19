package http

import (
	"context"
	"log"
	"net/http"
	"time"

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
}

func MustRun(ctx context.Context, shutdownDur time.Duration, addr string, app service) {
	handler := Handler{
		service: app,
	}

	r := mux.NewRouter()
	r.HandleFunc("/", handler.HealthCheck)
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)
	r.HandleFunc("/ticket/{ticket_id}", handler.GetTicketByIDHandler)
	r.HandleFunc("/ticket/check/{ticket_id}", handler.CheckTicketHandler)

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

}

func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello"))
}

func (h *Handler) CheckTicketHandler(w http.ResponseWriter, r *http.Request) {
	ticketID, err := types.GetTicketIDRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	valid, err := h.service.CheckTicket(r.Context(), ticketID)

	types.ProcessError(w, err, &types.CheckTicketHandlerResponse{Valid: valid})
}

func (h *Handler) GetTicketByIDHandler(w http.ResponseWriter, r *http.Request) {
	ticketID, err := types.GetTicketIDRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	ticket, err := h.service.GetTicket(r.Context(), ticketID)

	types.ProcessError(w, err, &types.GetTicketByIDHandlerResponse{Ticket: ticket})
}

type Handler struct {
	service service
}
