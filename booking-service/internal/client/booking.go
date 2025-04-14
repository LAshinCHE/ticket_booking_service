package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type SagaClient struct {
	Client http.Client
}

func NewSagaClient() *SagaClient {
	return &SagaClient{
		Client: http.Client{Timeout: time.Duration(1) * time.Second},
	}
}

func (sc *SagaClient) StartBookingSaga(ctx context.Context, bookingID int64, userID int64, ticketID int64) error {
	body := map[string]interface{}{
		"booking_id": bookingID,
		"user_id":    userID,
	}

	jsonData, _ := json.Marshal(body)
	resp, err := http.Post("http://saga-service:8080/start", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Saga service returned status %d", resp.StatusCode)
	}

	return nil
}
