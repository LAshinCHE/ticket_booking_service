package models

import "github.com/google/uuid"

type NotificationRequest struct {
	UserID  uuid.UUID
	Message string
}
