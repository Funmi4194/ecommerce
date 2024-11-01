package helper

import (
	"unicode"
)

// IsValidPassword checks if the provided password meets the expected requiremnet
func IsValidPassword(password string) (bool, error) {

	var hasSpecialChar, hasNumber, hasUpperCase, hasLowerCase bool
	for _, char := range password {
		if unicode.IsPunct(char) || unicode.IsSymbol(char) {
			hasSpecialChar = true
		} else if unicode.IsNumber(char) {
			hasNumber = true
		} else if unicode.IsUpper(char) {
			hasUpperCase = true
		} else if unicode.IsLetter(char) {
			hasLowerCase = true
		}
	}

	if !hasSpecialChar || !hasNumber || !hasLowerCase || !hasUpperCase {
		return false, nil
	}

	return true, nil
}
