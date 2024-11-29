package controller

import "github.com/labstack/echo/v4"

type OrderController interface {
	PlaceOrder(pctx echo.Context) error
	WorkerUpdateProgress(pctx echo.Context) error
	ConfirmationTimerOrder(pctx echo.Context) error
	CancelOrder(pctx echo.Context) error
	GetPreparingOrder(pctx echo.Context) error
	GetWorkerPrepare(pctx echo.Context) error
	GetQRpayment(pctx echo.Context) error
	InquiryQRpayment(pctx echo.Context) error
	GetWorkerFeePayment(pctx echo.Context) error
	InquiryFeePayment(pctx echo.Context) error

	GetUserOrder(pctx echo.Context) error
	GetUserContact(pctx echo.Context) error
	GetWorkerOrder(pctx echo.Context) error

	GetActiveOrder(pctx echo.Context) error
	GetActiveWorkerOrder(pctx echo.Context) error
}
