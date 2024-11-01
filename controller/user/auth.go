package user

import (
	"net/http"

	"github.com/funmi4194/ecommerce/helper"
	userLogic "github.com/funmi4194/ecommerce/logic/user"
	"github.com/funmi4194/ecommerce/types"
	"github.com/opensaucerer/barf"
)

// Register is the controller function to register a user
func Register(w http.ResponseWriter, r *http.Request) {

	var data types.User
	if err := barf.Request(r).Body().Format(&data); err != nil {
		barf.Logger().Errorf(`[user.Register] [barf.Request(r).Body().Format(&data)] %s`, err.Error())
		barf.Response(w).Status(http.StatusBadRequest).JSON(barf.Res{
			Status:  false,
			Message: "We could not process your request at this time. Please try again later.",
			Data:    nil,
		})
		return
	}

	user, err := userLogic.Register(data)
	if err != nil {
		barf.Logger().Errorf(`[user.Register] [userLogic.Register(data)] %s`, err.Error())
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
		Message: "Registration successful.",
		Data:    types.M{"user": user},
	})
}

// Login is the controller function to login a user
func Login(w http.ResponseWriter, r *http.Request) {

	var data types.Login
	if err := barf.Request(r).Body().Format(&data); err != nil {
		barf.Logger().Errorf(`[user.Login] [barf.Request(r).Body().Format(&data)] %s`, err.Error())
		barf.Response(w).Status(http.StatusBadRequest).JSON(barf.Res{
			Status:  false,
			Message: "We could not process your request at this time. Please try again later.",
			Data:    nil,
		})
		return
	}

	user, err := userLogic.Login(data)
	if err != nil {
		barf.Logger().Errorf(`[user.Login] [userLogic.Login(&data)] %s`, err.Error())
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
		Message: "Login successful.",
		Data: types.M{
			"user":  user,
			"token": helper.RefreshToken(user.ID),
		},
	})
}
