package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Token struct {
	ID         string
	Suffix     string
	Hash       string
	CreateTime time.Time
	ExpireTime time.Time
}

// NewToken creates a new token with the given expiration time.
// It returns the token and a pretty value that can be used to identify the token.
func NewToken(expireTime time.Time) (*Token, string) {
	value := generateRandomString(32)
	hash := tokenHash(value)

	prettyValue := fmt.Sprintf("sk_%s", value)
	suffix := prettyValue[len(prettyValue)-4:]

	return &Token{
		ID:         uuid.New().String(),
		Hash:       hash,
		Suffix:     suffix,
		CreateTime: time.Now(),
		ExpireTime: expireTime,
	}, prettyValue
}

func (t *Token) IsExpired() bool {
	return t.ExpireTime.Before(time.Now())
}

func tokenHash(token string) string {
	hasher := sha256.New()
	hasher.Write([]byte(token))
	hash := hasher.Sum(nil)

	return hex.EncodeToString(hash)
}
