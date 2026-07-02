package ids

import (
	"fmt"

	"github.com/google/uuid"
)

type SessionId struct{ Uuid }

func NewSessionId() SessionId { return SessionId{GenerateUuid()} }

func NewSessionIdFromString(raw string) (SessionId, error) {
	parsed, err := uuid.Parse(raw)
	if err != nil {
		return SessionId{}, fmt.Errorf("invalid session id: %w", err)
	}
	return SessionId{NewUuid(parsed)}, nil
}
