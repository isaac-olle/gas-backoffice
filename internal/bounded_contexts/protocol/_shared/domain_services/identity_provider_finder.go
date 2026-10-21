package domain_services

import (
	"context"

	protocolShared "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/_shared"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
)

type IdentityProviderFinder struct {
	repository protocolShared.IdentityProviderRepository
}

func NewIdentityProviderFinder(repository protocolShared.IdentityProviderRepository) *IdentityProviderFinder {
	return &IdentityProviderFinder{repository: repository}
}

func (f *IdentityProviderFinder) FindByReferenceName(ctx context.Context, referenceName string) (*protocolShared.IdentityProviderConfig, error) {
	idp, err := f.repository.FindByReferenceName(ctx, referenceName)
	if err != nil {
		return nil, app_errors.InternalError(protocolShared.ErrCannotGetIdentityProvider).
			WithDescription("repository failed to find identity provider").
			WithCause(err, false)
	}
	if idp == nil {
		return nil, app_errors.NotFoundError(protocolShared.ErrIdentityProviderNotFound).
			WithDescription("identity provider not found for reference name: " + referenceName)
	}
	return idp, nil
}
