package domain_services

import (
	"context"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/key_pairs"
)

type KeyPairFinder struct {
	repository key_pairs.Repository
}

func NewKeyPairFinder(repository key_pairs.Repository) *KeyPairFinder {
	return &KeyPairFinder{repository: repository}
}

func (f *KeyPairFinder) Find(ctx context.Context, id ids.KeyPairId) (*key_pairs.KeyPair, error) {
	kp, err := f.repository.SearchById(ctx, id)
	if err != nil {
		return nil, app_errors.InternalError(key_pairs.ErrKeyPairNotFound).
			WithDescription("repository failed to find key pair").
			WithCause(err)
	}
	if kp == nil {
		return nil, app_errors.NotFoundError(key_pairs.ErrKeyPairNotFound).
			WithDescription("key pair not found")
	}
	return kp, nil
}
