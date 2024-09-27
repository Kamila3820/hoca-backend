package controller

import "github.com/labstack/echo/v4"

type OrderController interface {
	PlaceOrder(pctx echo.Context) error
	WorkerUpdateProgress(pctx echo.Context) error
	ConfirmationTimerOrder(pctx echo.Context) error
	CancelOrder(pctx echo.Context) error
	GetPreparingOrder(pctx echo.Context) error

	GetUserContact(pctx echo.Context) error
}
