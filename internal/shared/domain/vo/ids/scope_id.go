package ids

import (
	"fmt"

	"github.com/google/uuid"
)

type ScopeId struct{ Uuid }

func NewScopeId() ScopeId { return ScopeId{GenerateUuid()} }

func NewScopeIdFromString(raw string) (ScopeId, error) {
	parsed, err := uuid.Parse(raw)
	if err != nil {
		return ScopeId{}, fmt.Errorf("invalid scope id: %w", err)
	}
	return ScopeId{NewUuid(parsed)}, nil
}
