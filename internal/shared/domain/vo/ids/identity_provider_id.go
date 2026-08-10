package ids

import (
	"fmt"

	"github.com/google/uuid"
)

type IdentityProviderId struct{ Uuid }

func NewIdentityProviderId() IdentityProviderId { return IdentityProviderId{GenerateUuid()} }

func NewIdentityProviderIdFromString(raw string) (IdentityProviderId, error) {
	parsed, err := uuid.Parse(raw)
	if err != nil {
		return IdentityProviderId{}, fmt.Errorf("invalid identity provider id: %w", err)
	}
	return IdentityProviderId{NewUuid(parsed)}, nil
}
