package primer

import (
	"crypto/sha256"
	"fmt"
)

// StringSha256 computes the sha256 of the given string
func StringSha256(input string) string {
	n := sha256.New()
	n.Write([]byte(input))
	return fmt.Sprintf("%x", n.Sum(nil))
}
