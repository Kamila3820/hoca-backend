package misc

import "github.com/golang-jwt/jwt/v5"

type UserClaim struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	jwt.RegisteredClaims
}

// GetAudience implements jwt.Claims.
func (T *UserClaim) GetAudience() (jwt.ClaimStrings, error) {
	panic("unimplemented")
}

// GetExpirationTime implements jwt.Claims.
func (T *UserClaim) GetExpirationTime() (*jwt.NumericDate, error) {
	return T.RegisteredClaims.ExpiresAt, nil
}

// GetIssuedAt implements jwt.Claims.
func (T *UserClaim) GetIssuedAt() (*jwt.NumericDate, error) {
	panic("unimplemented")
}

// GetIssuer implements jwt.Claims.
func (T *UserClaim) GetIssuer() (string, error) {
	panic("unimplemented")
}

// GetNotBefore implements jwt.Claims.
func (T *UserClaim) GetNotBefore() (*jwt.NumericDate, error) {
	return T.RegisteredClaims.NotBefore, nil
}

// GetSubject implements jwt.Claims.
func (T *UserClaim) GetSubject() (string, error) {
	panic("unimplemented")
}

func (T *UserClaim) Valid() error {
	return nil
}
