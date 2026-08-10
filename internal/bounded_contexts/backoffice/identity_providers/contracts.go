package identity_providers

import (
	"context"

	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

// Repository is read-only: identity providers are configured/seeded in the DB.
type Repository interface {
	FindById(ctx context.Context, id ids.IdentityProviderId) (*IdentityProvider, error)
}
