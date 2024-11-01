package helper

import (
	"github.com/opensaucerer/barf"
)

// RefreshToken sends token for every http request
func RefreshToken(id string) string {
	if id == "" {
		return ""
	}
	//generate new token
	token, err := SignJWT(id)
	if err != nil {
		barf.Logger().Errorf(`[helper.SendToken] [SignJWT(id)] %s`, err.Error())
		return ""
	}
	return token
}
