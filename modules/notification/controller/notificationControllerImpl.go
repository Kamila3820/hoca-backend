package controller

import (
	"net/http"

	"github.com/Kamila3820/hoca-backend/modules/account/misc"
	"github.com/Kamila3820/hoca-backend/modules/custom"
	_notificationService "github.com/Kamila3820/hoca-backend/modules/notification/service"
	"github.com/golang-jwt/jwt/v5"
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
	userID := pctx.Get("user").(*jwt.Token).Claims.(*misc.UserClaim)

	notifications, err := c.notificationService.GetNotificationsByUser(userID.ID)
	if err != nil {
		return custom.Error(pctx, http.StatusInternalServerError, err)
	}

	return pctx.JSON(http.StatusOK, notifications)
}
