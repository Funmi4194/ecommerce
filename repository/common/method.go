package common

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"encoding/json"

	"github.com/funmi4194/ecommerce/database"
	"github.com/uptrace/bun"
)

// BeginSaveTx returns a new transaction for the database
func BeginTx() (*bun.Tx, error) {
	tx, err := database.PostgreSQLDB.BeginTx(context.Background(), &sql.TxOptions{})
	return &tx, err
}

// Scan implements the Scanner interface.
func (h *History) Scan(src interface{}) error {
	switch v := src.(type) {
	case []byte:
		return json.Unmarshal(v, h)
	case string:
		return json.Unmarshal([]byte(v), h)
	case nil:
		return nil
	}
	return nil
}

// Scan implements the Scanner interface.
func (h *Histories) Scan(src interface{}) error {
	switch v := src.(type) {
	case []byte:
		return json.Unmarshal(v, h)
	case string:
		return json.Unmarshal([]byte(v), h)
	case nil:
		return nil
	}
	return nil
}

// Value implements the driver Valuer interface.
func (h Histories) Value() (driver.Value, error) {
	b, err := json.Marshal(h)
	return string(b), err
}

// Value implements the driver Valuer interface.
func (h History) Value() (driver.Value, error) {
	b, err := json.Marshal(h)
	return string(b), err
}
