package product

import (
	"net/http"

	"github.com/funmi4194/ecommerce/helper"
	"github.com/funmi4194/ecommerce/logic/product"
	"github.com/funmi4194/ecommerce/repository/user"
	"github.com/funmi4194/ecommerce/types"
	"github.com/opensaucerer/barf"
)

func Store(w http.ResponseWriter, r *http.Request) {

	// get user from context
	userId := r.Context().Value(types.AuthCtxKey{}).(*user.User).ID

	store, err := product.Store(barf.Request(r).Form().File().All("object"), userId)
	if err != nil {
		barf.Logger().Errorf(`[product.Store] [product.Store(barf.Request(r).Form().File().All("object"))] %s`, err.Error())
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
		Message: "Object(s) stored sucessfully",
		Data: types.M{
			"store": store,
			"token": helper.RefreshToken(userId),
		},
	})
}
