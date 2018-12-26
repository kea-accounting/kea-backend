package util

import (
	"encoding/base64"
	"math/rand"
)

// IsValidName returns whether a string is a valid name.   A String is a valid name if
// it consists entirely of [A-Z], [a-z], [0-9], '_', or '-'.
func IsValidName(text string) bool {

	for _, r := range text {
		if !isValidName(r) {
			return false
		}
	}

	return true
}

func isValidName(r rune) bool {
	return ((r >= 'A') && (r <= 'Z')) ||
		((r >= 'a') && (r <= 'z')) ||
		((r >= '0') && (r <= '9')) ||
		(r == '_') ||
		(r == '-')
}

// NewID generates a random 16-byte value, encoded as a URL-safe base64 string
func NewID() string {
	id := make([]byte, 16)
	rand.Read(id)
	return base64.RawURLEncoding.EncodeToString([]byte(id))
}
