package product

import (
	"github.com/funmi4194/ecommerce/controller/product"
	"github.com/opensaucerer/barf"
)

func RegisterStorageRoutes(frame *barf.SubRoute) {

	frame = frame.RetroFrame("/products")

	frame.Post("/store", product.Store)
}
