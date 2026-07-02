package ids

import (
	"fmt"

	"github.com/google/uuid"
)

type ProtectedResourceId struct{ Uuid }

func NewProtectedResourceId() ProtectedResourceId { return ProtectedResourceId{GenerateUuid()} }

func NewProtectedResourceIdFromString(raw string) (ProtectedResourceId, error) {
	parsed, err := uuid.Parse(raw)
	if err != nil {
		return ProtectedResourceId{}, fmt.Errorf("invalid protected resource id: %w", err)
	}
	return ProtectedResourceId{NewUuid(parsed)}, nil
}
