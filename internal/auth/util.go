package auth

import (
	"crypto/rand"
	"math/big"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// generateRandomString generates a cryptographically secure random string of the specified
// length using characters from charset.
func generateRandomString(length int) string {
	b := make([]rune, length)
	charsetLength := big.NewInt(int64(len(charset)))

	for i := range b {
		n, err := rand.Int(rand.Reader, charsetLength)
		if err != nil {
			panic("failed to generate random string: " + err.Error())
		}
		b[i] = []rune(charset)[n.Int64()]
	}
	return string(b)
}
