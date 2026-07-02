package ids

import (
	"fmt"

	"github.com/google/uuid"
)

type KeyPairId struct{ Uuid }

func NewKeyPairId() KeyPairId { return KeyPairId{GenerateUuid()} }

func NewKeyPairIdFromString(raw string) (KeyPairId, error) {
	parsed, err := uuid.Parse(raw)
	if err != nil {
		return KeyPairId{}, fmt.Errorf("invalid key pair id: %w", err)
	}
	return KeyPairId{NewUuid(parsed)}, nil
}
