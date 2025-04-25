package clients

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type PaymentClient struct {
	Address string
	Client  http.Client
}

func NewPaymentClient(addr string) *PaymentClient {
	return &PaymentClient{
		Address: addr,
		Client:  http.Client{Timeout: time.Duration(1) * time.Second},
	}
}

type ChargeRequest struct {
	UserID int64   `json:"user_id"`
	Amount float64 `json:"amount"`
}

type ChargeResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

func (c *PaymentClient) Charge(ctx context.Context, userID int64, amount float64) error {
	url := fmt.Sprintf("%s/payment/charge", c.Address)

	reqBody := ChargeRequest{
		UserID: userID,
		Amount: amount,
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal charge request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(bodyBytes))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.Client.Do(req)
	if err != nil {
		return fmt.Errorf("payment service request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("payment service returned status: %d", resp.StatusCode)
	}

	var chargeResp ChargeResponse
	if err := json.NewDecoder(resp.Body).Decode(&chargeResp); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	if !chargeResp.Success {
		return fmt.Errorf("payment failed: %s", chargeResp.Message)
	}

	return nil
}
