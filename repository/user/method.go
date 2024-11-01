package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/funmi4194/ecommerce/database"
	"github.com/funmi4194/ecommerce/enum"
	"github.com/funmi4194/ecommerce/primer"
	"github.com/funmi4194/ecommerce/reflection"
	"github.com/funmi4194/ecommerce/types"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/schema"
)

// Validate validates the user struct
func (u *User) Prepare() error {
	if u.Email == "" {
		return errors.New("email is required")
	}
	if u.Password == "" {
		return errors.New("password is required")
	}

	if len(u.Password) < primer.MinPassword {
		return fmt.Errorf("password must be greater than %d", primer.MinPassword-1)
	}

	u.Email = strings.ToLower(u.Email)

	return nil
}

/*
Date loads the created_at and updated_at fields of the activity if not already present, otherwise, it loads the updated_at field only.

If the "pessimistic" parameter is set to true, it loads both fields regardless
*/
func (u *User) Date(pessimistic ...bool) {
	if len(pessimistic) > 0 && !pessimistic[0] {
		if u.CreatedAt.IsZero() {
			u.CreatedAt = schema.NullTime{Time: time.Now()}
			u.UpdatedAt = schema.NullTime{Time: time.Now()}
			return
		}
		u.UpdatedAt = schema.NullTime{Time: time.Now()}
		return
	}
	u.CreatedAt = schema.NullTime{Time: time.Now()}
	u.UpdatedAt = schema.NullTime{Time: time.Now()}
}

/* Fields returns the struct fields as a slice of interface{} values */
func (u *User) Fields() []interface{} {
	return reflection.ReturnStructFields(u)
}

// Create inserts a new user into the database.
func (u *User) Create() error {
	if _, err := database.PostgreSQLDB.NewRaw(`INSERT INTO users (id, email, password, role, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)`, u.ID, u.Email, u.Password, u.Role, u.CreatedAt, u.UpdatedAt).Exec(context.Background()); err != nil && err != sql.ErrNoRows {
		return err
	}
	return nil
}

// CreateTx inserts a new user into the database with the given transaction
func (u *User) CreateTx(tx *bun.Tx) error {
	if _, err := tx.NewRaw(`INSERT INTO users (id, email, password, role, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)`, u.ID, u.Email, u.Password, u.Role, u.CreatedAt, u.UpdatedAt).Exec(context.Background()); err != nil && err != sql.ErrNoRows {
		return err
	}
	return nil
}

/*
FByKeyVal finds and returns a user matching the key/value pair

# By default, only the id and name fields are loaded

The	"preloadandjoin" parameter can be used to request that all the fields of the struct be loaded

It returns an error if any
*/
func (u *User) FByKeyVal(key string, val interface{}, preloadandjoin ...bool) error {
	if len(preloadandjoin) > 1 && preloadandjoin[0] && preloadandjoin[1] {
		return database.PostgreSQLDB.NewRaw(`SELECT * FROM users WHERE users.`+key+` = ?`, val).Scan(context.Background(), u.Fields()...)
	}
	if len(preloadandjoin) > 0 && preloadandjoin[0] {
		return database.PostgreSQLDB.NewRaw(`SELECT * FROM users WHERE `+key+` = ?`, val).Scan(context.Background(), u)
	}
	return database.PostgreSQLDB.NewRaw(`SELECT id, email FROM users WHERE users.`+key+` = ?`, val).Scan(context.Background(), u)
}

func (u *User) Execute(q string, args ...interface{}) error {
	return database.PostgreSQLDB.NewRaw(q, args...).Scan(context.Background(), u)
}

/*
UByMap updates a user matching the key/value pairs provided in the map

It returns an error if any
*/
func (u *User) UByMap(m types.SQLMaps) error {
	query, args := database.MapsToSQuery(m)
	if strings.Contains(query, string(enum.RETURNING)) {
		return database.PostgreSQLDB.NewRaw(`UPDATE users `+query, args...).Scan(context.Background(), u)
	}
	_, err := database.PostgreSQLDB.NewRaw(`UPDATE users `+query, args...).Exec(context.Background())
	return err
}
