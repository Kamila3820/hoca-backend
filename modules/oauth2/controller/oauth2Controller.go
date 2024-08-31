package controller

import "github.com/labstack/echo/v4"

type OAuth2Controller interface {
	UserLogin(pctx echo.Context) error
	UserLoginCallback(pctx echo.Context) error
	Logout(pctx echo.Context) error

	UserAuthorizing(pctx echo.Context, next echo.HandlerFunc) error
}
