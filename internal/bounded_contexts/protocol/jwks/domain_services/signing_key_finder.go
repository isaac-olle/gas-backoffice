package domain_services

import (
	"context"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/jwks"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
)

type SigningKeyFinder struct {
	repository jwks.SigningKeyRepository
}

func NewSigningKeyFinder(repository jwks.SigningKeyRepository) *SigningKeyFinder {
	return &SigningKeyFinder{repository: repository}
}

func (f *SigningKeyFinder) FindActive(ctx context.Context) ([]*jwks.SigningKeyProjection, error) {
	keys, err := f.repository.SearchActive(ctx)
	if err != nil {
		return nil, app_errors.InternalError(jwks.ErrCannotGetSigningKeys).
			WithDescription("repository failed to find signing keys").
			WithCause(err, false)
	}
	return keys, nil
}
