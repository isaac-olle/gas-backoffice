package ids

import (
	"fmt"

	"github.com/google/uuid"
)

type AdminAssertionKeyId struct{ Uuid }

func NewAdminAssertionKeyId() AdminAssertionKeyId {
	return AdminAssertionKeyId{GenerateUuid()}
}

func NewAdminAssertionKeyIdFromString(raw string) (AdminAssertionKeyId, error) {
	parsed, err := uuid.Parse(raw)
	if err != nil {
		return AdminAssertionKeyId{}, fmt.Errorf("invalid admin assertion key id: %w", err)
	}
	return AdminAssertionKeyId{NewUuid(parsed)}, nil
}
