package controller

import "github.com/labstack/echo/v4"

type NotificationController interface {
	GetNotificationsByUser(pctx echo.Context) error
}
