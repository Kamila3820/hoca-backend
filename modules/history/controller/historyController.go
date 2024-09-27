package controller

import "github.com/labstack/echo/v4"

type HistoryController interface {
	GetOrderHistoryByUserID(pctx echo.Context) error
	GetHistoryByUserID(pctx echo.Context) error
	GetWorkingHistory(pctx echo.Context) error
}
