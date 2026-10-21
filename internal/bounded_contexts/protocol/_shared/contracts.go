package shared

import "context"

type IdentityProviderRepository interface {
	FindByReferenceName(ctx context.Context, referenceName string) (*IdentityProviderConfig, error)
}
