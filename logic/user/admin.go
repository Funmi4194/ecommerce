package user

import (
	"database/sql"
	"errors"

	"github.com/funmi4194/ecommerce/enum"
	userRepository "github.com/funmi4194/ecommerce/repository/user"
	"github.com/funmi4194/ecommerce/types"
	"github.com/opensaucerer/barf"
)

/*
AddAdmin is a function that would allow an admin user add a user as an admin
this will help ensure testing admin functionanlities
*/
func AddAdmin(userId string, payload types.AdminPayload) error {

	user := userRepository.User{
		ID: userId,
	}

	// find user by ID
	err := user.FByKeyVal("id", user.ID, true)
	if err != nil {
		barf.Logger().Errorf(`[user.AddAdmin] [user.FByKeyVal("id", user.ID, true)] %s`, err.Error())
		if err == sql.ErrNoRows {
			return errors.New("looks like your account no longer exists. please contact support")
		}
		return errors.New("we're having issues retrieving your account. please try again later")
	}

	if user.Role != enum.Admin {
		return errors.New("you do not have the permission to add an admin")
	}

	// find user by ID
	err = user.FByKeyVal("id", payload.UserID, true)
	if err != nil {
		barf.Logger().Errorf(`[user.AddAdmin] [user.FByKeyVal("id", user.ID, true)] %s`, err.Error())
		if err == sql.ErrNoRows {
			return errors.New("looks like your account no longer exists. please contact support")
		}
		return errors.New("we're having issues adding user as admin. please try again later")
	}

	if user.Role == enum.Admin {
		return errors.New("user already have access to admin feature")
	}

	// add user to admin
	err = user.UByMap(types.SQLMaps{
		WMaps: []types.SQLMap{
			{
				Map: map[string]interface{}{
					"id": user.ID,
				},
				ComparisonOperator: enum.Equal,
			},
		},
		SMap: types.SQLMap{
			Map: map[string]interface{}{
				"role":       enum.Admin,
				"updated_at": "now()",
			},
			ComparisonOperator: enum.Equal,
			JoinOperator:       enum.Comma,
		},
	})
	if err != nil {
		barf.Logger().Errorf(`[user.AddAdmin] [user.UByMap(query)] %s`, err.Error())
		if err == sql.ErrNoRows {
			return errors.New("looks like your account no longer exists. please contact support")
		}
		return errors.New("we're having issues adding user as admin. please try again later")
	}

	return nil
}
