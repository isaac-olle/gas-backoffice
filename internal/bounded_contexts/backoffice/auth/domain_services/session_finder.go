package domain_services

import (
	"context"

	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/auth"
)

type SessionFinder struct {
	repository auth.Repository
}

func NewSessionFinder(repository auth.Repository) *SessionFinder {
	return &SessionFinder{repository: repository}
}

func (f *SessionFinder) Find(ctx context.Context, id ids.SessionId) (*auth.Session, error) {
	session, err := f.repository.SearchById(ctx, id)
	if err != nil {
		return nil, app_errors.InternalError(auth.ErrSessionNotFound).
			WithDescription("error searching session").
			WithCause(err)
	}
	if session == nil {
		return nil, app_errors.NotFoundError(auth.ErrSessionNotFound).
			WithDescription("session not found")
	}
	return session, nil
}
