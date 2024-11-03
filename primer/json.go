package primer

import (
	"encoding/json"
)

// Stringify is a helper function to convert an interface to a string.
func Stringify(i interface{}) string {
	b, _ := json.Marshal(i)
	return string(b)
}
