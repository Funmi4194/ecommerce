package user

import (
	"net/http"

	"github.com/funmi4194/ecommerce/helper"
	userLogic "github.com/funmi4194/ecommerce/logic/user"
	"github.com/funmi4194/ecommerce/repository/user"
	"github.com/funmi4194/ecommerce/types"
	"github.com/opensaucerer/barf"
)

// AddAdmin is the controller function to add an admin user
func AddAdmin(w http.ResponseWriter, r *http.Request) {

	// get user from context
	id := r.Context().Value(types.AuthCtxKey{}).(*user.User).ID

	var data types.AdminPayload
	if err := barf.Request(r).Body().Format(&data); err != nil {
		barf.Logger().Errorf(`[user.AddAdmin] [barf.Request(r).Body().Format(&data)] %s`, err.Error())
		barf.Response(w).Status(http.StatusBadRequest).JSON(barf.Res{
			Status:  false,
			Message: "We could not process your request at this time. Please try again later.",
			Data:    nil,
		})
		return
	}

	err := userLogic.AddAdmin(id, data)
	if err != nil {
		barf.Logger().Errorf(`[user.AddAdmin] [userLogic.AddAdmin(id, data)] %s`, err.Error())
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
		Message: "Admin added successfully",
		Data: types.M{
			"token": helper.RefreshToken(id),
		},
	})
}
