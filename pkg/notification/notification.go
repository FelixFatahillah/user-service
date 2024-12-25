package notification

import (
	"context"
	"encoding/json"
	"user-service/internal/infrastructure/mq/publisher"
)

type NotificationType string

const (
	OrderSuccess NotificationType = "order_success"
	OrderExpired NotificationType = "order_expired"
)

type NotificationSender string

const (
	EmailSender    NotificationSender = "email"
	WhatsappSender NotificationSender = "whatsapp"
)

type SendNotification struct {
	Message            json.RawMessage    `json:"message" validate:"required"`
	Subject            string             `json:"subject"`
	NotificationType   NotificationType   `json:"notification_type" validate:"required"`
	Target             string             `json:"target" validate:"required"`
	NotificationSender NotificationSender `json:"notification_sender" validate:"required"`
}

func Sender(ctx context.Context, data SendNotification) error {
	message := map[string]interface{}{
		"message":             data.Message,
		"subject":             data.Subject,
		"service_source":      "ms-order",
		"notification_type":   data.NotificationType,
		"target":              data.Target,
		"notification_sender": data.NotificationSender,
	}
	jsonMessage, err := json.Marshal(message)
	body := publisher.Publish{
		Headers: nil,
		Body:    jsonMessage,
	}
	err = publisher.SendNotificationRoute.Publish(ctx, &body)
	if err != nil {
		return err
	}
	return nil
}
