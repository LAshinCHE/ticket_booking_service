package clients

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/LAshinCHE/ticket_booking_service/saga-orchestrator/models"
)

type NotificationClient struct {
	Addres     string
	HTTPClient http.Client
}

func NewNotificationClient(addr string, timeout time.Duration) *NotificationClient {
	return &NotificationClient{
		Addres: addr,
		HTTPClient: http.Client{
			Timeout: timeout,
		},
	}
}

func (nc *NotificationClient) Notify(ctx context.Context, params models.BookingParams) error {
	url := fmt.Sprintf("http//localhost:%d/notify/%d", nc.Addres, params.UserID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := nc.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to call booking service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("booking check failed with status: %d", resp.StatusCode)
	}

	return nil
}
