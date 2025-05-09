package http

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/LAshinCHE/ticket_booking_service/notification-service/internal/models"
	"github.com/gorilla/mux"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

type NotificationService interface {
	SendNotification(ctx context.Context, req models.NotificationRequest) error
}

func MustRun(ctx context.Context, addr string, app NotificationService, shutdowmDur time.Duration) {
	h := Handler{
		service: app,
	}

	r := mux.NewRouter()
	r.HandleFunc("/", h.HealthCheck).Methods("GET")
	r.HandleFunc("/notify", h.Notify).Methods("GET")

	server := http.Server{
		Addr:    addr,
		Handler: r,
	}
	go func() {
		<-ctx.Done()

		log.Printf("Shuting down server with duration %0.3fs", shutdowmDur.Seconds())

		<-time.After(shutdowmDur)
		if err := server.Shutdown(context.Background()); err != nil {
			log.Fatalf("Http handler Shutdown: %s", err)
		}
	}()

	if err := server.ListenAndServe(); err != nil {
		log.Printf("HTTP server ListenAndServe: %s", err)
	}
}

func (h *Handler) Notify(w http.ResponseWriter, r *http.Request) {
	logPrefix := "[NotifyService]"
	propagator := propagation.TraceContext{}
	ctx := propagator.Extract(r.Context(), propagation.HeaderCarrier(r.Header))

	tracer := otel.Tracer("notification-service")
	ctx, span := tracer.Start(ctx, "Notify")
	defer span.End()

	type NotifyRequest struct {
		UserID  int64  `json:"user_id"`
		Message string `json:"message"`
	}

	var reqData NotifyRequest
	if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
		log.Printf("%s ошибка при чтении тела запроса: %v", logPrefix, err)
		http.Error(w, "неверный формат запроса", http.StatusBadRequest)
		return
	}

	h.service.SendNotification(
		r.Context(),
		models.NotificationRequest{
			UserID:  reqData.UserID,
			Message: reqData.Message,
		})
}

func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello"))
}

type Handler struct {
	service NotificationService
}
