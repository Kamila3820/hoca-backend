package controller

import "github.com/labstack/echo/v4"

type OrderController interface {
	PlaceOrder(pctx echo.Context) error
}
