package controller

import (
	"context"
	"net/http"

	"github.com/Kamila3820/hoca-backend/modules/custom"
	_oauth2Exception "github.com/Kamila3820/hoca-backend/modules/oauth2/exception"
	"github.com/labstack/echo/v4"
	"golang.org/x/oauth2"
)

func (c *googleOAuth2Controller) UserAuthorizing(pctx echo.Context, next echo.HandlerFunc) error {
	ctx := context.Background()

	tokenSource, err := c.getTokenSource(pctx)
	if err != nil {
		return custom.Error(pctx, http.StatusUnauthorized, err)
	}

	if !tokenSource.Valid() {
		tokenSource, err = c.userTokenRefreshing(pctx, tokenSource)
		if err != nil {
			return custom.Error(pctx, http.StatusUnauthorized, err)
		}
	}

	client := userGoogleOAuth2.Client(ctx, tokenSource)

	userInfo, err := c.getUserInfo(client)
	if err != nil {
		return custom.Error(pctx, http.StatusUnauthorized, err)
	}

	if !c.oauth2Service.IsUserExists(userInfo.ID) {
		return custom.Error(pctx, http.StatusUnauthorized, &_oauth2Exception.Unauthorized{})
	}

	pctx.Set("userID", userInfo.ID)

	return next(pctx)
}

func (c *googleOAuth2Controller) userTokenRefreshing(pctx echo.Context, token *oauth2.Token) (*oauth2.Token, error) {
	ctx := context.Background()

	updatedToken, err := userGoogleOAuth2.TokenSource(ctx, token).Token()
	if err != nil {
		return nil, &_oauth2Exception.Unauthorized{}
	}

	c.setSameSiteCookie(pctx, accessTokenCookieName, updatedToken.AccessToken)
	c.setSameSiteCookie(pctx, refreshTokenCookieName, updatedToken.RefreshToken)

	return updatedToken, nil
}

func (c *googleOAuth2Controller) getTokenSource(pctx echo.Context) (*oauth2.Token, error) {
	accessToken, err := pctx.Cookie(accessTokenCookieName)
	if err != nil {
		return nil, &_oauth2Exception.Unauthorized{}
	}

	refreshToken, err := pctx.Cookie(refreshTokenCookieName)
	if err != nil {
		return nil, &_oauth2Exception.Unauthorized{}
	}

	return &oauth2.Token{
		AccessToken:  accessToken.Value,
		RefreshToken: refreshToken.Value,
	}, nil
}
