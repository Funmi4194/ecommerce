package helper

import (
	"fmt"
	"strings"
	"time"

	"github.com/funmi4194/ecommerce/primer"
	"github.com/google/uuid"
)

// GenerateUUID returns a new universally unique identifier
func GenerateUUID() string {
	return uuid.New().String()
}

// GenerateFilename generates a filename from a given filename
func GenerateFilename(filename string) string {
	return fmt.Sprintf(`%s-%s`, GenerateUUID(), filename)
}

// GenerateRef returns a usable reference computed against the given label (or app name) from the current timestamp
func GenerateRef(label ...string) string {
	if len(label) == 0 {
		label = []string{primer.ENV.AppName.String()}
	}
	return fmt.Sprintf(`%s-%d`, strings.ToLower(label[0]), time.Now().Nanosecond())
}
