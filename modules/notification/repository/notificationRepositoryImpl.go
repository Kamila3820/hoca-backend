package repository

import (
	"github.com/Kamila3820/hoca-backend/entities"
	"github.com/Kamila3820/hoca-backend/pkg/databases"
	"github.com/labstack/echo/v4"
)

type notificationRepositoryImpl struct {
	db     databases.Database
	logger echo.Logger
}

func NewNotificationRepositoryImpl(db databases.Database, logger echo.Logger) NotificationRepository {
	return &notificationRepositoryImpl{
		db:     db,
		logger: logger,
	}
}

func (r *notificationRepositoryImpl) GetNotificationByUser(userID string) ([]*entities.Notification, error) {
	var notifications []*entities.Notification

	db := r.db.Connect().Preload("Trigger").Preload("Triggee").Preload("Order").Preload("UserRating").Order("created_at DESC")

	if err := db.Where("triggee_id = ?", userID).Find(&notifications).Error; err != nil {
		r.logger.Errorf("Unable to query notifications: %s", err.Error())
		return nil, err
	}

	return notifications, nil
}
