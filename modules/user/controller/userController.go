package controller

import "github.com/labstack/echo/v4"

type UserController interface {
	GetUserByID(pctx echo.Context) error
	EditUserProfile(pctx echo.Context) error
	GetUserLocation(pctx echo.Context) error
	UpdateUserLocation(pctx echo.Context) error
}
