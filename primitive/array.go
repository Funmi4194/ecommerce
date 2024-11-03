package primitive

import (
	"database/sql/driver"
	"encoding/json"
	"strings"
)

// I have a plan to make a package from this called nanites.
type StringArray []string

// ExistsIn reports whether any of the elements of sa is contained in t.
func (sa StringArray) ExistsIn(t string) bool {
	for _, v := range sa {
		if strings.Contains(t, v) {
			return true
		}
	}
	return false
}

// Scan implements the Scanner interface.
func (sa *StringArray) Scan(src interface{}) error {
	switch v := src.(type) {
	case []byte:
		return json.Unmarshal(v, sa)
	case string:
		return json.Unmarshal([]byte(v), sa)
	case nil:
		return nil
	}
	return nil
}

// Value implements the driver Valuer interface.
func (sa StringArray) Value() (driver.Value, error) {
	b, err := json.Marshal(sa)
	return string(b), err
}

// Find returns the first element in sa that satisfies the provided testing function.
// Otherwise nil is returned.
func (sa StringArray) Find(fn func(interface{}) bool) interface{} {
	for _, v := range sa {
		if fn(v) {
			return v
		}
	}
	return nil
}

type Array []interface{}
