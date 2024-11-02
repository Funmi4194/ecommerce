package version

import (
	"github.com/funmi4194/ecommerce/middleware"
	"github.com/funmi4194/ecommerce/route/product"
	"github.com/funmi4194/ecommerce/route/user"
	"github.com/opensaucerer/barf"
)

func V1() {
	unauthenticedFrame := barf.RetroFrame("/v1")
	barf.Hippocampus(unauthenticedFrame).Hijack(middleware.OptionalAuth)

	// access to some parts of the api is only allowed with a valid token
	authenticatedFrame := barf.RetroFrame("/v1")
	barf.Hippocampus(authenticatedFrame).Hijack(middleware.Auth)

	user.RegisterAuthRoutes(unauthenticedFrame)
	user.RegisterAdminRoutes(authenticatedFrame)

	product.RegisterProductRoutes(authenticatedFrame)
}
