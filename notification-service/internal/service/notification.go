package service

import (
	"context"
	"log"

	"github.com/LAshinCHE/ticket_booking_service/notification-service/internal/models"
)

type NotificationService struct {
}

func NewNotificationService() *NotificationService {
	return &NotificationService{}
}

func (s *NotificationService) SendNotification(ctx context.Context, req models.NotificationRequest) error {
	log.Printf("Notify user %s: %s\n", req.UserID, req.Message)
	// Здесь может быть логика отправки email, push и т.п.
	return nil
}
