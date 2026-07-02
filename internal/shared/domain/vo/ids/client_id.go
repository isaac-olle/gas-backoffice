package ids

import (
	"fmt"

	"github.com/google/uuid"
)

type ClientId struct{ Uuid }

func NewClientId() ClientId { return ClientId{GenerateUuid()} }

func NewClientIdFromString(raw string) (ClientId, error) {
	parsed, err := uuid.Parse(raw)
	if err != nil {
		return ClientId{}, fmt.Errorf("invalid client id: %w", err)
	}
	return ClientId{NewUuid(parsed)}, nil
}
