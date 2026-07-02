package ids

import (
	"fmt"

	"github.com/google/uuid"
)

type TenantId struct{ Uuid }

func NewTenantId() TenantId { return TenantId{GenerateUuid()} }

func NewTenantIdFromString(raw string) (TenantId, error) {
	parsed, err := uuid.Parse(raw)
	if err != nil {
		return TenantId{}, fmt.Errorf("invalid tenant id: %w", err)
	}
	return TenantId{NewUuid(parsed)}, nil
}
