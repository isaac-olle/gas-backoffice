package domain_services

import (
	"context"
	"github.com/isaac-olle/gas/internal/shared/app_errors"
	ids2 "github.com/isaac-olle/gas/internal/shared/domain/vo/ids"

	"github.com/isaac-olle/gas/internal/bounded_contexts/security/client_assertion_keys"
)

type ClientAssertionKeyRepository interface {
	client_assertion_keys.Repository
	FindByClientId(ctx context.Context, clientID ids2.ClientId) (*client_assertion_keys.ClientAssertionKey, error)
}

type ClientAssertionKeyFinder struct {
	repository ClientAssertionKeyRepository
}

func NewClientAssertionKeyFinder(repository ClientAssertionKeyRepository) *ClientAssertionKeyFinder {
	return &ClientAssertionKeyFinder{repository: repository}
}

func (f *ClientAssertionKeyFinder) Find(ctx context.Context, id ids2.ClientAssertionKeyId) (*client_assertion_keys.ClientAssertionKey, error) {
	cak, err := f.repository.SearchById(ctx, id)
	if err != nil {
		return nil, err
	}
	if cak == nil {
		return nil, app_errors.NotFoundError(client_assertion_keys.ErrClientAssertionKeyNotFound).
			WithDescription("client assertion key not found")
	}
	return cak, nil
}

func (f *ClientAssertionKeyFinder) FindByClientId(ctx context.Context, clientID ids2.ClientId) (*client_assertion_keys.ClientAssertionKey, error) {
	cak, err := f.repository.FindByClientId(ctx, clientID)
	if err != nil {
		return nil, err
	}
	if cak == nil {
		return nil, app_errors.NotFoundError(client_assertion_keys.ErrClientAssertionKeyNotFound).
			WithDescription("client assertion key not found")
	}
	return cak, nil
}
