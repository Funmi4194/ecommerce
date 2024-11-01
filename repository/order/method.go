package order

import (
	"database/sql/driver"
	"encoding/json"
)

// Scan implements the Scanner interface.
func (i *Item) Scan(src interface{}) error {
	switch v := src.(type) {
	case []byte:
		return json.Unmarshal(v, i)
	case string:
		return json.Unmarshal([]byte(v), i)
	case nil:
		return nil
	}
	return nil
}

// Value implements the driver Valuer interface.
func (i Item) Value() (driver.Value, error) {
	b, err := json.Marshal(i)
	return string(b), err
}
