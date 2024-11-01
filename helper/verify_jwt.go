package helper

import (
	"github.com/funmi4194/ecommerce/primer"
	"github.com/funmi4194/ecommerce/types"
	"github.com/golang-jwt/jwt/v4"
)

// VerifyJWT verifies a JWT and returns the claims
func VerifyJWT(token string) (*types.JWTClaims, bool) {
	if twc, err := jwt.ParseWithClaims(token, &types.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(primer.ENV.JWTSecret), nil
	}); err == nil && twc.Valid {
		return twc.Claims.(*types.JWTClaims), true
	}
	return nil, false
}
