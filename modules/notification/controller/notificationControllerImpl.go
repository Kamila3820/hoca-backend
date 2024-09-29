package controller

import (
	"net/http"

	"github.com/Kamila3820/hoca-backend/modules/custom"
	_notificationService "github.com/Kamila3820/hoca-backend/modules/notification/service"
	"github.com/labstack/echo/v4"
)

type notificationControllerImpl struct {
	notificationService _notificationService.NotificationService
}

func NewNotificationControllerImpl(notificationService _notificationService.NotificationService) NotificationController {
	return &notificationControllerImpl{
		notificationService: notificationService,
	}
}

func (c *notificationControllerImpl) GetNotificationsByUser(pctx echo.Context) error {
	userID := pctx.Get("userID")
	userIDStr, ok := userID.(string)
	if !ok {
		return pctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to retrieve user ID from context",
		})
	}

	notifications, err := c.notificationService.GetNotificationsByUser(userIDStr)
	if err != nil {
		return custom.Error(pctx, http.StatusInternalServerError, err)
	}

	return pctx.JSON(http.StatusOK, notifications)
}
