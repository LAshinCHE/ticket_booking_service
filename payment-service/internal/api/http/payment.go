package http

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type service interface {
	// CheckBalance(ctx context.Context, userID uuid.UUID, amount float64) (bool, error)
	DebitBalance(ctx context.Context, userID int64, amount float64) (bool, error)
	RefundBalance(ctx context.Context, userID int64, amount float64) error
}

func MustRun(ctx context.Context, app service, addr string, shutdownDur time.Duration) {
	h := Handler{
		service: app,
	}

	r := mux.NewRouter()
	r.HandleFunc("/health-check/", h.HealthCheck).Methods("GET")
	r.HandleFunc("/payments/charge", h.DebitFromBalance).Methods("POST")
	r.HandleFunc("/payments/refund", h.RefundToBalance).Methods("POST")
	// r.HandleFunc("/payments/{user_id}/balance/", h.GetBalance).Methods("GET")
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

// Debit также проверяет хватает ли средств на счету
func (h *Handler) DebitFromBalance(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIDStr := vars["userID"]
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid userID", http.StatusBadRequest)
		return
	}

	amountStr := r.URL.Query().Get("amount")
	if amountStr == "" {
		http.Error(w, "missing amount", http.StatusBadRequest)
		return
	}

	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		http.Error(w, "invalid amount", http.StatusBadRequest)
		return
	}

	enough, err := h.service.DebitBalance(r.Context(), userID, amount)
	if err != nil {
		if err.Error() == "insufficient funds" {
			http.Error(w, "not enough funds", http.StatusConflict)
			return
		}
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	resp := struct {
		Enough bool `json:"enough"`
	}{
		Enough: enough,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func (h *Handler) RefundToBalance(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIDStr := vars["userID"]
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid userID", http.StatusBadRequest)
		return
	}

	amountStr := r.URL.Query().Get("amount")
	if amountStr == "" {
		http.Error(w, "missing amount", http.StatusBadRequest)
		return
	}

	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		http.Error(w, "invalid amount", http.StatusBadRequest)
		return
	}

	err = h.service.RefundBalance(r.Context(), userID, amount)

	if err != nil {
		if err.Error() == "insufficient funds" {
			http.Error(w, "not enough funds", http.StatusConflict)
			return
		}
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

type Handler struct {
	service service
}
