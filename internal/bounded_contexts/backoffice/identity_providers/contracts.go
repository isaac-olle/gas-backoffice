package identity_providers

import (
	"context"

	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

type Repository interface {
	Save(ctx context.Context, idp *IdentityProvider) error
	FindById(ctx context.Context, id ids.IdentityProviderId) (*IdentityProvider, error)
	FindByReferenceName(ctx context.Context, referenceName string) (*IdentityProvider, error)
}
