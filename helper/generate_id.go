package helper

import (
	"github.com/google/uuid"
)

// GenerateUUID returns a new universally unique identifier
func GenerateUUID() string {
	return uuid.New().String()
}
