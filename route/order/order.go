package order

import (
	orderController "github.com/funmi4194/ecommerce/controller/order"
	"github.com/opensaucerer/barf"
)

func RegisterOrderRoutes(frame *barf.SubRoute) {

	frame = frame.RetroFrame("/order")

	frame.Post("/create", orderController.InitiateOrder)
	frame.Post("/list", orderController.Orders)
	frame.Patch("/update", orderController.UpdateOrder)
	frame.Patch("/cancel", orderController.CancelOrder)
}
