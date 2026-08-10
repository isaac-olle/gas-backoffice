package domain_services

import (
	"context"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
)

type IdentityProviderFinder struct {
	repository session.IdentityProviderFinderRepository
}

func NewIdentityProviderFinder(repository session.IdentityProviderFinderRepository) *IdentityProviderFinder {
	return &IdentityProviderFinder{repository: repository}
}

func (f *IdentityProviderFinder) Find(ctx context.Context, idpID string) (*session.IdentityProvider, error) {
	idp, err := f.repository.FindById(ctx, idpID)
	if err != nil {
		return nil, app_errors.InternalError(session.ErrCannotGetIdentityProvider).
			WithDescription("repository failed to find identity provider").
			WithCause(err)
	}
	if idp == nil {
		return nil, app_errors.NotFoundError(session.ErrIdentityProviderNotFound).
			WithDescription("identity provider not found for id: " + idpID)
	}
	return idp, nil
}
