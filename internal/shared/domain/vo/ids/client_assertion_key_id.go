package ids

import (
	"fmt"

	"github.com/google/uuid"
)

type ClientAssertionKeyId struct{ Uuid }

func NewClientAssertionKeyId() ClientAssertionKeyId {
	return ClientAssertionKeyId{GenerateUuid()}
}

func NewClientAssertionKeyIdFromString(raw string) (ClientAssertionKeyId, error) {
	parsed, err := uuid.Parse(raw)
	if err != nil {
		return ClientAssertionKeyId{}, fmt.Errorf("invalid client assertion key id: %w", err)
	}
	return ClientAssertionKeyId{NewUuid(parsed)}, nil
}
