package http

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type service interface{}

func MustRun(ctx context.Context, app service, addr string, shutdownDur time.Duration) {
	h := Handler{
		service: app,
	}

	r := mux.NewRouter()
	r.HandleFunc("/health-check/", h.HealthCheck).Methods("GET")
	// прописать route

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

func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello"))
}

type Handler struct {
	service service
}
