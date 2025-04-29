package http

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/LAshinCHE/ticket_booking_service/notification-service/internal/models"
	"github.com/gorilla/mux"
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
	message := r.URL.Query().Get("message")
	if message == "" {
		http.Error(w, "missing message", http.StatusBadRequest)
		return
	}

	userIDStr := r.URL.Query().Get("userID")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		http.Error(w, "missing message", http.StatusBadRequest)
		return
	}

	h.service.SendNotification(
		r.Context(),
		models.NotificationRequest{
			UserID:  userID,
			Message: message,
		})
}

func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello"))
}

type Handler struct {
	service NotificationService
}
