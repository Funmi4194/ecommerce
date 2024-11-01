package helper

import (
	"time"

	"github.com/funmi4194/ecommerce/primer"
	"github.com/funmi4194/ecommerce/types"
	"github.com/golang-jwt/jwt/v4"
)

// SignJWT signs a JWT with the given address
func SignJWT(id string, durations ...time.Duration) (string, error) {
	var expr time.Duration

	// check if a duration was provided, if not use the default duration
	if len(durations) > 0 {
		expr = durations[0]
	} else {
		// set default duration
		expr = 24 * 3 * time.Hour
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, types.JWTClaims{
		ID: id,
		RegisteredClaims: jwt.RegisteredClaims{
			// expires in 24 hours
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expr)),
			Issuer:    primer.ENV.AppName.String(),
		},
	}).SignedString([]byte(primer.ENV.JWTSecret))
}
