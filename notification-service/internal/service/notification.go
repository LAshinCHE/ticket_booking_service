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
	log.Printf("Notify user %d: %s\n", req.UserID, req.Message)
	return nil
}
