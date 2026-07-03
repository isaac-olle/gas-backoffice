package domain_services

import (
	"context"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session"
)

type ClientFinder struct {
	repository session.ClientFinderRepository
}

func NewClientFinder(repository session.ClientFinderRepository) *ClientFinder {
	return &ClientFinder{repository: repository}
}

func (f *ClientFinder) Find(ctx context.Context, clientID string) (*session.Client, error) {
	client, err := f.repository.FindById(ctx, clientID)
	if err != nil {
		return nil, app_errors.InternalError(session.ErrCannotGetClient).
			WithDescription("repository failed to find client").
			WithCause(err)
	}
	if client == nil {
		return nil, app_errors.NotFoundError(session.ErrClientNotFound).
			WithDescription("client not found for clientID: " + clientID)
	}
	return client, nil
}
