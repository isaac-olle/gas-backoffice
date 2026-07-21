package ids

import (
	"fmt"

	"github.com/google/uuid"
)

type UserAssertionKeyId struct{ Uuid }

func NewUserAssertionKeyId() UserAssertionKeyId {
	return UserAssertionKeyId{GenerateUuid()}
}

func NewUserAssertionKeyIdFromString(raw string) (UserAssertionKeyId, error) {
	parsed, err := uuid.Parse(raw)
	if err != nil {
		return UserAssertionKeyId{}, fmt.Errorf("invalid user assertion key id: %w", err)
	}
	return UserAssertionKeyId{NewUuid(parsed)}, nil
}
