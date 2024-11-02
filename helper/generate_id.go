package helper

import (
	"fmt"

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
