package service

import (
	_notiModel "github.com/Kamila3820/hoca-backend/modules/notification/model"
)

type NotificationService interface {
	GetNotificationsByUser(userID string) ([]*_notiModel.NotificationResponse, error)
}
