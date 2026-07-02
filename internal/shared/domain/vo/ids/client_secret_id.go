package ids

import (
	"fmt"

	"github.com/google/uuid"
)

type ClientSecretId struct{ Uuid }

func NewClientSecretId() ClientSecretId { return ClientSecretId{GenerateUuid()} }

func NewClientSecretIdFromString(raw string) (ClientSecretId, error) {
	parsed, err := uuid.Parse(raw)
	if err != nil {
		return ClientSecretId{}, fmt.Errorf("invalid client secret id: %w", err)
	}
	return ClientSecretId{NewUuid(parsed)}, nil
}
