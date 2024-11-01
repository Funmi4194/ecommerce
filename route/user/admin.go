package user

import (
	userController "github.com/funmi4194/ecommerce/controller/user"
	"github.com/opensaucerer/barf"
)

func RegisterAdminRoutes(frame *barf.SubRoute) {

	frame = frame.RetroFrame("/accounts")

	frame.Post("/add/admin", userController.AddAdmin)
}
