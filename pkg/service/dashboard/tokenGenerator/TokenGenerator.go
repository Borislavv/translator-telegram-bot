package tokenGenerator

import (
	"crypto/rand"
	"encoding/hex"
)

var (
	salt        int = 9
	tokenLength int = 118
)

type TokenGenerator struct {
}

// NewTokenGenerator - constructor of SessionTokenGenerator struct
func NewTokenGenerator() *TokenGenerator {
	return &TokenGenerator{}
}

// Generate - generate a session token
func (generator *TokenGenerator) Generate() (string, error) {
	bytes := make([]byte, tokenLength)
	bytes = append(bytes, make([]byte, salt)...)

	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	return hex.EncodeToString(bytes), nil
}
