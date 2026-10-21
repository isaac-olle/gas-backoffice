package domain_service

import (
	"context"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

type ClientProjectionFinder struct {
	repository clients.Repository
}

func NewClientFinder(repository clients.Repository) *ClientProjectionFinder {
	return &ClientProjectionFinder{repository: repository}
}

func (f *ClientProjectionFinder) Find(ctx context.Context, id ids.ClientId) (*clients.Client, error) {
	client, err := f.repository.SearchById(ctx, id)
	if err != nil {
		return nil, app_errors.InternalError(clients.ErrClientNotFound).
			WithDescription("failed to search client").WithCause(

			err, false)

	}
	if client == nil {
		return nil, app_errors.NotFoundError(clients.ErrClientNotFound).
			WithDescription("client not found")
	}
	return client, nil
}
