package domain_services

import (
	"context"
	"github.com/isaac-olle/gas/internal/shared/app_errors"
	"github.com/isaac-olle/gas/internal/shared/domain/vo/ids"

	"github.com/isaac-olle/gas/internal/bounded_contexts/backoffice/clients"
)

type ClientFinder struct {
	repository clients.Repository
}

func NewClientFinder(repository clients.Repository) *ClientFinder {
	return &ClientFinder{repository: repository}
}

func (f *ClientFinder) Find(ctx context.Context, id ids.ClientId) (*clients.Client, error) {
	client, err := f.repository.SearchById(ctx, id)
	if err != nil {
		return nil, app_errors.InternalError(clients.ErrClientNotFound).
			WithDescription("failed to search client").
			WithCause(err)
	}
	if client == nil {
		return nil, app_errors.NotFoundError(clients.ErrClientNotFound).
			WithDescription("client not found")
	}
	return client, nil
}
