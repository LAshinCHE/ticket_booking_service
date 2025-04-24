package clients

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type BookingClient struct {
	BaseURL    string
	HTTPClient *http.Client
}

func NewBookingClient(baseURL string) *BookingClient {
	return &BookingClient{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (c *BookingClient) CheckBooking(ctx context.Context, bookingID int64) error {
	url := fmt.Sprintf("%s/booking/%d/check", c.BaseURL, bookingID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to call booking service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("booking check failed with status: %d", resp.StatusCode)
	}

	return nil
}
