package client_secrets

import (
	"context"
	ids2 "github.com/isaac-olle/gas/internal/shared/domain/vo/ids"
	"github.com/isaac-olle/gas/internal/shared/persistance"
)

type Repository interface {
	persistance.Repository[*ClientSecret, ids2.ClientSecretId]
	FindByClientId(ctx context.Context, clientID ids2.ClientId) (*ClientSecret, error)
}

type SecretHasher interface {
	Hash(plainSecret string) (string, error)
}

type SecretVerifier interface {
	Verify(ctx context.Context, clientID, plainSecret string) error
}
