package product

import (
	"net/http"

	"github.com/funmi4194/ecommerce/helper"
	"github.com/funmi4194/ecommerce/logic/product"
	"github.com/funmi4194/ecommerce/repository/user"
	"github.com/funmi4194/ecommerce/types"
	"github.com/opensaucerer/barf"
)

// Publish is the controller function to create product or products
func Publish(w http.ResponseWriter, r *http.Request) {

	// get user from context
	userId := r.Context().Value(types.AuthCtxKey{}).(*user.User).ID

	var data types.Publish
	if err := barf.Request(r).Body().Format(&data); err != nil {
		barf.Logger().Errorf(`[product.Publish] [barf.Request(r).Body().Format(&data)] %s`, err.Error())
		barf.Response(w).Status(http.StatusBadRequest).JSON(barf.Res{
			Status:  false,
			Message: "We could not process your request at this time. Please try again later.",
			Data:    nil,
		})
		return
	}

	products, err := product.Publish(userId, data)
	if err != nil {
		barf.Logger().Errorf(`[product.Publish] [product.Publish(data, userId)] %s`, err.Error())
		barf.Response(w).Status(http.StatusBadRequest).JSON(barf.Res{
			Status:  false,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	// send response
	barf.Response(w).Status(http.StatusCreated).JSON(barf.Res{
		Status:  true,
		Message: "Product(s) published sucessfully",
		Data: types.M{
			"products": products,
			"token":    helper.RefreshToken(userId),
		},
	})
}

// UpdateProduct is the controller function to update a product
func UpdateProduct(w http.ResponseWriter, r *http.Request) {

	// get user from context
	userId := r.Context().Value(types.AuthCtxKey{}).(*user.User).ID

	var data types.UpdateProduct
	if err := barf.Request(r).Body().Format(&data); err != nil {
		barf.Logger().Errorf(`[product.UpdateProduct] [barf.Request(r).Body().Format(&data)] %s`, err.Error())
		barf.Response(w).Status(http.StatusBadRequest).JSON(barf.Res{
			Status:  false,
			Message: "We could not process your request at this time. Please try again later.",
			Data:    nil,
		})
		return
	}

	product, err := product.UpdateProduct(userId, data)
	if err != nil {
		barf.Logger().Errorf(`[product.UpdateProduct] [product.UpdateProduct(userId, data)] %s`, err.Error())
		barf.Response(w).Status(http.StatusBadRequest).JSON(barf.Res{
			Status:  false,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	// send response
	barf.Response(w).Status(http.StatusOK).JSON(barf.Res{
		Status:  true,
		Message: "Product(s) retreived sucessfully",
		Data: types.M{
			"product": product,
			"token":   helper.RefreshToken(userId),
		},
	})
}

// Products is the controller function to fetch products
func Products(w http.ResponseWriter, r *http.Request) {

	// get user from context
	userId := r.Context().Value(types.AuthCtxKey{}).(*user.User).ID

	var data types.ProductFilter
	if err := barf.Request(r).Body().Format(&data); err != nil {
		barf.Logger().Errorf(`[product.Product] [barf.Request(r).Body().Format(&data)] %s`, err.Error())
		barf.Response(w).Status(http.StatusBadRequest).JSON(barf.Res{
			Status:  false,
			Message: "We could not process your request at this time. Please try again later.",
			Data:    nil,
		})
		return
	}

	products, pagination, err := product.Products(userId, data)
	if err != nil {
		barf.Logger().Errorf(`[product.Product] [product.Products(userId, data)] %s`, err.Error())
		barf.Response(w).Status(http.StatusBadRequest).JSON(barf.Res{
			Status:  false,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	// send response
	barf.Response(w).Status(http.StatusOK).JSON(barf.Res{
		Status:  true,
		Message: "Product(s) retreived sucessfully",
		Data: types.M{
			"products":   products,
			"pagination": pagination,
			"token":      helper.RefreshToken(userId),
		},
	})
}

// Product is the controller function to fetch a product by its Id
func Product(w http.ResponseWriter, r *http.Request) {

	userId := r.Context().Value(types.AuthCtxKey{}).(*user.User).ID

	var data types.ProductFilter
	if err := barf.Request(r).Query().Format(&data); err != nil {
		barf.Logger().Errorf(`[product.Product] [barf.Request(r).Body().Format(&data)] %s`, err.Error())
		barf.Response(w).Status(http.StatusBadRequest).JSON(barf.Res{
			Status:  false,
			Message: "We could not process your request at this time. Please try again later.",
			Data:    nil,
		})
		return
	}

	product, err := product.Product(userId, data)
	if err != nil {
		barf.Logger().Errorf(`[product.Product] product.Product(userId, data)] %s`, err.Error())
		barf.Response(w).Status(http.StatusBadRequest).JSON(barf.Res{
			Status:  false,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	// send response
	barf.Response(w).Status(http.StatusOK).JSON(barf.Res{
		Status:  true,
		Message: "Product retreived sucessfully",
		Data: types.M{
			"product": product,
			"token":   helper.RefreshToken(userId),
		},
	})
}
