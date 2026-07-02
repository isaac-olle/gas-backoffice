package infrastructure

import (
	"fmt"

	"github.com/alexedwards/argon2id"
)

type Argon2SecretHasher struct{}

func NewArgon2SecretHasher() *Argon2SecretHasher { return &Argon2SecretHasher{} }

func (h *Argon2SecretHasher) Hash(plainSecret string) (string, error) {
	hash, err := argon2id.CreateHash(plainSecret, argon2id.DefaultParams)
	if err != nil {
		return "", fmt.Errorf("failed to hash secret: %w", err)
	}
	return hash, nil
}
