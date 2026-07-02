package ids

import (
	"fmt"

	"github.com/google/uuid"
)

type ClientCredentialsSessionId struct{ Uuid }

func NewClientCredentialsSessionId() ClientCredentialsSessionId {
	return ClientCredentialsSessionId{GenerateUuid()}
}

func NewClientCredentialsSessionIdFromString(raw string) (ClientCredentialsSessionId, error) {
	parsed, err := uuid.Parse(raw)
	if err != nil {
		return ClientCredentialsSessionId{}, fmt.Errorf("invalid client credentials session id: %w", err)
	}
	return ClientCredentialsSessionId{NewUuid(parsed)}, nil
}
