package primitive

import (
	"database/sql/driver"
	"encoding/json"
	"strings"
)

// I have a plan to make a package from this called nanites. Basically, nanites will expose multi primitives implemented in the Go way. https://github.com/opensaucerer/nanites
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

type Array []interface{}
