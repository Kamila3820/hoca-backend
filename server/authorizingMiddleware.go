package server

import (
	"github.com/Kamila3820/hoca-backend/config"
	_oauth2Controller "github.com/Kamila3820/hoca-backend/modules/oauth2/controller"
	"github.com/labstack/echo/v4"
)

type authorizingMiddleware struct {
	oauth2Controller _oauth2Controller.OAuth2Controller
	oauth2Conf       *config.OAuth2
	logger           echo.Logger
}

func (m *authorizingMiddleware) UserAuthorizing(next echo.HandlerFunc) echo.HandlerFunc {
	return func(pctx echo.Context) error {
		return m.oauth2Controller.UserAuthorizing(pctx, next)
	}
}
