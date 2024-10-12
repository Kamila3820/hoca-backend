package controller

import "github.com/labstack/echo/v4"

type AccountController interface {
	Register(pctx echo.Context) error
	Login(pctx echo.Context) error
}
