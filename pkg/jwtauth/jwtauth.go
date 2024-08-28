package jwtauth

import "github.com/golang-jwt/jwt/v5"

type (
	AuthFactory interface {
		SignToken() string
	}

	Claims struct {
		PlayerId string `json:"player_id"`
		RoleCode int    `json:"role_code"`
	}

	AuthMapClaims struct {
		*Claims
		jwt.RegisteredClaims
	}

	authConcrete struct {
		Secret []byte
		Claims *AuthMapClaims `json:"claims"`
	}

	accessToken  struct{ *authConcrete }
	refreshToken struct{ *authConcrete }
)
