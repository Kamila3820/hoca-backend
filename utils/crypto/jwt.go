package crypto

import (
	"errors"

	"github.com/golang-jwt/jwt/v4"
)

func SignJwt(claim jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	if signedToken, err := token.SignedString([]byte("babycomeandtakemylovenadruinit")); err != nil {
		return "", errors.New("Unable to sign JWT token")
	} else {
		return signedToken, nil
	}
}
