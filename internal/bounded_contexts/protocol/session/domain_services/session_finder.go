package domain_services

import (
	"context"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session"
)

type SessionFinder struct {
	repository session.Repository
}

func NewSessionFinder(repository session.Repository) *SessionFinder {
	return &SessionFinder{repository: repository}
}

func (f *SessionFinder) FindByCode(ctx context.Context, code string) (*session.AuthorizationCodeSession, error) {
	sess, err := f.repository.FindByCode(ctx, code)
	if err != nil {
		return nil, app_errors.InternalError(session.ErrCodeNotFound).
			WithDescription("repository failed to find session by code").
			WithCause(err)
	}
	if sess == nil {
		return nil, app_errors.BadRequestError(session.ErrCodeNotFound).
			WithDescription("authorization code not found or expired")
	}
	return sess, nil
}

func (f *SessionFinder) FindById(ctx context.Context, id ids.SessionId) (*session.AuthorizationCodeSession, error) {
	sess, err := f.repository.FindById(ctx, id)
	if err != nil {
		return nil, app_errors.InternalError(session.ErrSessionNotFound).
			WithDescription("repository failed to find session").
			WithCause(err)
	}
	if sess == nil {
		return nil, app_errors.NotFoundError(session.ErrSessionNotFound).
			WithDescription("session not found")
	}
	return sess, nil
}
