package ids

import (
	"fmt"

	"github.com/google/uuid"
)

type AdminId struct{ Uuid }

func NewAdminId() AdminId { return AdminId{GenerateUuid()} }

func NewAdminIdFromString(raw string) (AdminId, error) {
	parsed, err := uuid.Parse(raw)
	if err != nil {
		return AdminId{}, fmt.Errorf("invalid admin id: %w", err)
	}
	return AdminId{NewUuid(parsed)}, nil
}
