package key_pairs

import (
	"context"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"github.com/isaac-olle/gas-backoffice/internal/shared/persistance"
	"time"
)

type Repository interface {
	persistance.Repository[*KeyPair, ids.KeyPairId]
	SearchByProtectedResourceId(ctx context.Context, id ids.ProtectedResourceId) (*KeyPair, error)
}

type TokenVerifier interface {
	Verify(ctx context.Context, token string, kp *KeyPair, expectedClaims map[string]any, verifyExpiration bool) (map[string]any, error)
}

type TokenGenerator interface {
	Generate(ctx context.Context, iss string, spec TokenSpec, kp *KeyPair, expiresIn time.Duration) (string, error)
}

type JwtDecoder interface {
	ExtractKeyId(ctx context.Context, token string) (string, error)
}

type KeyGenerator interface {
	Generate(algorithm string) (privateKeyPEM string, publicKeyPEM string, err error)
}
