package user

import (
	userController "github.com/funmi4194/ecommerce/controller/user"
	"github.com/opensaucerer/barf"
)

func RegisterAuthRoutes(frame *barf.SubRoute) {

	frame = frame.RetroFrame("/accounts")

	frame.Post("/register", userController.Register)
	frame.Post("/login", userController.Login)
}
