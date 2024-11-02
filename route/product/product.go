package product

import (
	productController "github.com/funmi4194/ecommerce/controller/product"
	"github.com/opensaucerer/barf"
)

func RegisterProductRoutes(frame *barf.SubRoute) {

	frame = frame.RetroFrame("/products")

	frame.Post("/publish", productController.Publish)
	frame.Patch("/update", productController.UpdateProduct)
	frame.Get("/product", productController.Product)
	frame.Delete("/delete", productController.DeleteProduct)
	frame.Post("/list", productController.Products)
}
