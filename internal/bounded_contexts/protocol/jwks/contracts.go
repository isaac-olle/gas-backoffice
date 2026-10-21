package jwks

import "context"

type SigningKeyRepository interface {
	SearchActive(ctx context.Context) ([]*SigningKeyProjection, error)
}
