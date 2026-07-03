package key_pairs

import (
	"context"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"github.com/isaac-olle/gas-backoffice/internal/shared/persistance"
	"time"
)

type Repository interface {
	persistance.Repository[*KeyPair, ids.KeyPairId]
}

type TokenVerifier interface {
	Verify(ctx context.Context, token string, kp *KeyPair) (*TokenClaims, error)
}

type TokenGenerator interface {
	Generate(ctx context.Context, claims map[string]any, kp *KeyPair, expiresIn time.Duration) (string, error)
}

type JwtDecoder interface {
	ExtractKeyId(ctx context.Context, token string) (string, error)
}
