package repository

import "github.com/Kamila3820/hoca-backend/entities"

type NotificationRepository interface {
	GetNotificationByUser(userID string) ([]*entities.Notification, error)
}
