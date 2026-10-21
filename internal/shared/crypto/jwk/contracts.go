package jwk

import (
	"context"

	sharedCrypto "github.com/isaac-olle/gas-backoffice/internal/shared/crypto"
)

// KeyFetcher fetches the VerifyingKeys published at a JWKS endpoint.
type KeyFetcher interface {
	FetchAll(ctx context.Context, jwksUri string) (map[string]sharedCrypto.VerifyingKey, error)
	FetchByKid(ctx context.Context, jwksUri, kid string) (sharedCrypto.VerifyingKey, error)
}
