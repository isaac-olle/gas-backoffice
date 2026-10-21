package domain_services

import (
	"context"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

type ClientCredentialsSessionFinder struct {
	repository session.ClientCredentialsSessionRepository
}

func NewClientCredentialsSessionFinder(repository session.ClientCredentialsSessionRepository) *ClientCredentialsSessionFinder {
	return &ClientCredentialsSessionFinder{repository: repository}
}

func (f *ClientCredentialsSessionFinder) FindById(ctx context.Context, id ids.SessionId) (*session.ClientCredentialsSession, error) {
	sess, err := f.repository.FindById(ctx, id)
	if err != nil {
		return nil, app_errors.InternalError(session.ErrSessionNotFound).
			WithDescription("repository failed to find client credentials session").WithCause(err, false)
	}
	if sess == nil {
		return nil, app_errors.NotFoundError(session.ErrSessionNotFound).
			WithDescription("client credentials session not found")
	}
	return sess, nil
}
