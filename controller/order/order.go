package order

import (
	"net/http"

	"github.com/funmi4194/ecommerce/helper"
	"github.com/funmi4194/ecommerce/logic/order"
	"github.com/funmi4194/ecommerce/repository/user"
	"github.com/funmi4194/ecommerce/types"
	"github.com/opensaucerer/barf"
)

// InitiateOrder is the controller function to create order
func InitiateOrder(w http.ResponseWriter, r *http.Request) {

	// get user from context
	userId := r.Context().Value(types.AuthCtxKey{}).(*user.User).ID

	var data types.InitiateOrder
	if err := barf.Request(r).Body().Format(&data); err != nil {
		barf.Logger().Errorf(`[order.InitiateOrder] [barf.Request(r).Body().Format(&data)] %s`, err.Error())
		barf.Response(w).Status(http.StatusBadRequest).JSON(barf.Res{
			Status:  false,
			Message: "We could not process your request at this time. Please try again later.",
			Data:    nil,
		})
		return
	}

	order, err := order.InitiateOrder(userId, data)
	if err != nil {
		barf.Logger().Errorf(`[order.InitiateOrder] [order.InitiateOrder(userId, data)] %s`, err.Error())
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
		Message: "Order(s) initiated sucessfully",
		Data: types.M{
			"order": order,
			"token": helper.RefreshToken(userId),
		},
	})
}

// CancelOrder is the controller function to cancel order
func CancelOrder(w http.ResponseWriter, r *http.Request) {

	// get user from context
	userId := r.Context().Value(types.AuthCtxKey{}).(*user.User).ID

	var data types.CancelOrder
	if err := barf.Request(r).Body().Format(&data); err != nil {
		barf.Logger().Errorf(`[order.CancelOrder] [barf.Request(r).Body().Format(&data)] %s`, err.Error())
		barf.Response(w).Status(http.StatusBadRequest).JSON(barf.Res{
			Status:  false,
			Message: "We could not process your request at this time. Please try again later.",
			Data:    nil,
		})
		return
	}

	order, err := order.CancelOrder(userId, data)
	if err != nil {
		barf.Logger().Errorf(`[order.CancelOrder] [order.CancelOrder(userId, data)] %s`, err.Error())
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
		Message: "Order(s) cancelled sucessfully",
		Data: types.M{
			"order": order,
			"token": helper.RefreshToken(userId),
		},
	})
}

// UpdateOrder is the controller function to update an  order
func UpdateOrder(w http.ResponseWriter, r *http.Request) {

	// get user from context
	userId := r.Context().Value(types.AuthCtxKey{}).(*user.User).ID

	var data types.UpdateOrder
	if err := barf.Request(r).Body().Format(&data); err != nil {
		barf.Logger().Errorf(`[order.UpdateOrder] [barf.Request(r).Body().Format(&data)] %s`, err.Error())
		barf.Response(w).Status(http.StatusBadRequest).JSON(barf.Res{
			Status:  false,
			Message: "We could not process your request at this time. Please try again later.",
			Data:    nil,
		})
		return
	}

	order, err := order.UpdateOrder(userId, data)
	if err != nil {
		barf.Logger().Errorf(`[order.UpdateOrder] [order.UpdateOrder(userId, data)] %s`, err.Error())
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
		Message: "Order(s) updated sucessfully",
		Data: types.M{
			"order": order,
			"token": helper.RefreshToken(userId),
		},
	})
}

// Orders is the controller function to retrieve orders
func Orders(w http.ResponseWriter, r *http.Request) {

	// get user from context
	userId := r.Context().Value(types.AuthCtxKey{}).(*user.User).ID

	var data types.OrderFilter
	if err := barf.Request(r).Body().Format(&data); err != nil {
		barf.Logger().Errorf(`[order.Orders] [barf.Request(r).Body().Format(&data)] %s`, err.Error())
		barf.Response(w).Status(http.StatusBadRequest).JSON(barf.Res{
			Status:  false,
			Message: "We could not process your request at this time. Please try again later.",
			Data:    nil,
		})
		return
	}

	order, pagination, err := order.Orders(userId, data)
	if err != nil {
		barf.Logger().Errorf(`[order.Orders] [order.Orders(userId, data)] %s`, err.Error())
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
		Message: "Order(s) retrieved sucessfully",
		Data: types.M{
			"order":      order,
			"pagination": pagination,
			"token":      helper.RefreshToken(userId),
		},
	})
}
