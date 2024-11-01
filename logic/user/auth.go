package user

import (
	"database/sql"
	"errors"
	"net/mail"
	"strings"

	"github.com/funmi4194/ecommerce/enum"
	"github.com/funmi4194/ecommerce/helper"
	"github.com/funmi4194/ecommerce/primer"
	commonRepository "github.com/funmi4194/ecommerce/repository/common"
	userRepository "github.com/funmi4194/ecommerce/repository/user"
	"github.com/funmi4194/ecommerce/types"
	"github.com/opensaucerer/barf"
	"golang.org/x/crypto/bcrypt"
)

// Register registers a new user
func Register(payload types.User) (*userRepository.User, error) {

	user := userRepository.User{
		Email:    payload.Email,
		Password: payload.Password,
	}

	// verify user payload
	if err := user.Prepare(); err != nil {
		barf.Logger().Errorf(`[user.Register] [user.Prepare()] %s`, err.Error())
		return nil, errors.New("something's not right. Please try again or open a support ticket")
	}

	// validate the user email address
	// this serves as an additional check to verify the email is a valid email format
	email, err := mail.ParseAddress(payload.Email)
	if err != nil {
		barf.Logger().Errorf(`[user.EmailVerification] [mail.ParseAddress(payload.Email)] %s`, err.Error())
		return nil, errors.New("there seem to be an issue with the email address you provided. Please provide a valid email address")
	}

	// ensure email is unique
	err = user.FByKeyVal("email", strings.ToLower(email.Address))
	if err == nil {
		return nil, errors.New("email address is already in use")
	}
	if err != sql.ErrNoRows {
		barf.Logger().Errorf(`[user.Register] [user.FByKeyVal("email", strings.ToLower(email.Address))] %s`, err.Error())
		return nil, errors.New("we are having issues verifying this email address. Please try again later")
	}

	// verify password
	isValid, err := helper.IsValidPassword(user.Password)
	if err != nil {
		barf.Logger().Errorf(`[user.Register] [helper.IsValidPassword(user.Password)] %s`, err.Error())
		return nil, errors.New("something's not right. Please try again or open a support ticket")
	}

	if !isValid {
		return nil, errors.New("password must include at least one uppercase letter, one lowercase letter, one special character, and one number")
	}

	// hash password
	hashed, err := bcrypt.GenerateFromPassword([]byte(payload.Password), primer.HashCost)
	if err != nil {
		barf.Logger().Errorf(`[user.UpdatePassword] [bcrypt.GenerateFromPassword([]byte(payload.Password), primer.HashCost)] %s`, err.Error())
		return nil, errors.New("we are having issues creating your account. Please try again later")
	}

	user.Password = string(hashed)
	user.Role = enum.User
	user.ID = helper.GenerateUUID()
	user.Date()

	// create a new transaction
	btx, err := commonRepository.BeginTx()
	if err != nil {
		return nil, err
	}
	defer btx.Rollback()

	// create user
	if err := user.CreateTx(btx); err != nil {
		barf.Logger().Errorf(`[user.Register] [user.CreateTx(btx)] %s`, err.Error())
		return nil, errors.New("we are having issues creating your account. Please try again later")
	}

	// commit transaction
	if err := btx.Commit(); err != nil {
		barf.Logger().Errorf(`[user.Register] [btx.Commit()] %s`, err.Error())
		return nil, errors.New("we are having issues creating your account. Please try again later")
	}

	user.Password = ""

	return &user, nil
}

// Login sign in a user
func Login(payload types.Login) (*userRepository.User, error) {

	if payload.Email == "" {
		return nil, errors.New("email is required to login")
	}
	if payload.Password == "" {
		return nil, errors.New("password is required to login")
	}

	payload.Email = strings.ToLower(payload.Email)

	user := userRepository.User{}

	err := user.FByKeyVal("email", payload.Email, true)
	if err != nil {
		if err == sql.ErrNoRows {
			barf.Logger().Errorf(`[user.Login] [user.FByMap(query)] %s`, err.Error())
			return &user, errors.New("invalid login credentials provided")
		}
		return nil, errors.New("we are having issues signing you in. Please try again later")
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password)); err != nil {
		barf.Logger().Errorf(`[user.Register] [bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password))] %s`, err.Error())
		return nil, errors.New("invalid login credentials provided")
	}

	user.Password = ""

	return &user, nil
}
