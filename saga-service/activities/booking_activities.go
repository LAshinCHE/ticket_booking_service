package activities

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

// --------- инфраструктура ----------------------------------------------------

// ServiceClients содержит базовые URL-ы ваших микросервисов + http-клиент.
type ServiceClients struct {
	HTTPClient *http.Client
	BookingURL, TicketURL,
	PaymentURL, NotifyURL string
}

// NewServiceClients удобен для DI в tests.
func NewServiceClients(booking, ticket, payment, notify string) ServiceClients {
	return ServiceClients{
		HTTPClient: &http.Client{Timeout: 5 * time.Second},
		BookingURL: booking,
		TicketURL:  ticket,
		PaymentURL: payment,
		NotifyURL:  notify,
	}
}

// BookingActivities группирует все операции саги.
type BookingActivities struct {
	SVC ServiceClients
}

// NewBookingActivities создаём в main() и регистрируем  ➜  worker.RegisterActivity(...)
func NewBookingActivities(svc ServiceClients) *BookingActivities { return &BookingActivities{SVC: svc} }

// 1. Создать бронирование
func (a *BookingActivities) CreateBooking(
	ctx context.Context,
	bookingID int,
	userID int,
	ticketID int,
	price float64,
	traceCtx map[string]string,
) (string, error) {
	propagator := propagation.TraceContext{}
	parentCtx := propagator.Extract(context.Background(), propagation.MapCarrier(traceCtx))

	tracer := otel.Tracer("saga-activities")
	ctx, span := tracer.Start(parentCtx, "Activity.CreateBooking")
	defer span.End()
	log.Println("Booking ID:", bookingID)
	payload := struct {
		BookingID int     `json:"booking_id"`
		UserID    int     `json:"user_id"`
		TicketID  int     `json:"ticket_id"`
		Price     float64 `json:"price"`
	}{bookingID, userID, ticketID, price}

	raw, _ := json.Marshal(payload)
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, a.SVC.BookingURL+"/internal/booking/create", bytes.NewReader(raw))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "super-secure-saga-token")
	propagator.Inject(ctx, propagation.HeaderCarrier(req.Header))

	resp, err := a.SVC.HTTPClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("http request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return "", fmt.Errorf("create booking failed: %d", resp.StatusCode)
	}

	var respData struct {
		BookingID string `json:"booking_id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&respData); err != nil {
		return "", fmt.Errorf("decode booking response: %w", err)
	}

	return respData.BookingID, nil
}

// 2. Проверка билета в наличии
func (a *BookingActivities) CheckTicketAvailability(
	ctx context.Context,
	ticketID int,
) (bool, error) {
	ticketStrID := strconv.Itoa(ticketID)
	resBody, err := a.doGET(ctx, a.SVC.TicketURL+"/ticket/"+ticketStrID+"/check")
	if err != nil {
		return false, err
	}

	var resp struct {
		Available bool `json:"available"`
	}
	if err := json.Unmarshal(resBody, &resp); err != nil {
		return false, fmt.Errorf("decode availability: %w", err)
	}
	return resp.Available, nil
}

// 3. Перевод билета в статус «забронирован»
func (a *BookingActivities) ReserveTicket(
	ctx context.Context,
	ticketID int,
) error {

	payload := struct {
		Status string `json:"status"`
	}{"reserved"}
	ticketStrId := strconv.Itoa(ticketID)
	_, err := a.doPUT(ctx, a.SVC.TicketURL+"/tickets/"+ticketStrId, payload)
	return err
}

// 4. Списать средства
// Сharge будет проверять что у нас хватает денег на балансе и в случае нехватки возвращать ошибку
func (a *BookingActivities) WithdrawMoney(
	ctx context.Context,
	userID int,
	amount float64,
) error {

	payload := struct {
		UserID int     `json:"user_id"`
		Amount float64 `json:"amount"`
	}{userID, amount}

	_, err := a.doPOST(ctx, a.SVC.PaymentURL+"/payments/charge", payload)
	return err
}

// 4*. Отмена операции списания средства
// Сharge будет проверять что у нас хватает денег на балансе и в случае нехватки возвращать ошибку
func (a *BookingActivities) CancelWithdrawMoney(
	ctx context.Context,
	userID int,
	amount float64,
) error {

	payload := struct {
		UserID int     `json:"user_id"`
		Amount float64 `json:"amount"`
	}{userID, amount}

	_, err := a.doPOST(ctx, a.SVC.PaymentURL+"/payments/refund", payload)
	return err
}

// 5. Уведомить пользователя
func (a *BookingActivities) NotifyUser(
	ctx context.Context,
	userID int,
	message string,
) error {

	payload := struct {
		UserID  int    `json:"user_id"`
		Message string `json:"message"`
	}{userID, message}

	_, err := a.doPOST(ctx, a.SVC.NotifyURL+"/notifications", payload)
	return err
}

// --------- * Компенсация — отмена бронирования -------------------------------

func (a *BookingActivities) CancelBooking(
	ctx context.Context,
	bookingID int,
	traceCtx map[string]string,
) error {
	propagator := propagation.TraceContext{}
	parentCtx := propagator.Extract(context.Background(), propagation.MapCarrier(traceCtx))

	tracer := otel.Tracer("saga-activities")
	ctx, span := tracer.Start(parentCtx, "Activity.CancelBooking")
	log.Println("SPAN TRACE ID:", span.SpanContext().TraceID().String())
	defer span.End()

	payload := struct {
		BookingID int `json:"booking_id"`
	}{bookingID}

	raw, _ := json.Marshal(payload)
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, a.SVC.BookingURL+"/internal/booking/delete", bytes.NewReader(raw))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "super-secure-saga-token")
	propagator.Inject(ctx, propagation.HeaderCarrier(req.Header))

	resp, err := a.SVC.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("cancel booking: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("cancel booking: unexpected status %d", resp.StatusCode)
	}
	return nil
}

// --------- вспомогательные HTTP-helpers --------------------------------------

func (a *BookingActivities) doPOST(ctx context.Context, url string, payload any) ([]byte, error) {
	return a.doWithBody(ctx, http.MethodPost, url, payload)
}

func (a *BookingActivities) doPUT(ctx context.Context, url string, payload any) ([]byte, error) {
	return a.doWithBody(ctx, http.MethodPut, url, payload)
}

func (a *BookingActivities) doGET(ctx context.Context, url string) ([]byte, error) {
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	return a.do(req)
}

func (a *BookingActivities) doWithBody(ctx context.Context, method, url string, payload any) ([]byte, error) {
	raw, _ := json.Marshal(payload)
	req, _ := http.NewRequestWithContext(ctx, method, url, bytes.NewReader(raw))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "super-secure-saga-token")
	return a.do(req)
}

func (a *BookingActivities) do(req *http.Request) ([]byte, error) {
	resp, err := a.SVC.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http request %s %s: %w", req.Method, req.URL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("%s %s returned status %d", req.Method, req.URL, resp.StatusCode)
	}
	return io.ReadAll(resp.Body)
}
