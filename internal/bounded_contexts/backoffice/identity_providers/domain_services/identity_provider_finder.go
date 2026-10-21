package domain_services

import (
	"context"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/identity_providers"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

type IdentityProviderFinder struct {
	repository identity_providers.Repository
}

func NewIdentityProviderFinder(repository identity_providers.Repository) *IdentityProviderFinder {
	return &IdentityProviderFinder{repository: repository}
}

// FindEnabledByReferenceName returns the provider for the given reference name, only if it exists and is enabled.
func (f *IdentityProviderFinder) FindEnabledByReferenceName(ctx context.Context, referenceName string) (*identity_providers.IdentityProvider, error) {
	idp, err := f.repository.FindByReferenceName(ctx, referenceName)
	if err != nil {
		return nil, app_errors.InternalError(identity_providers.ErrIdentityProviderNotFound).
			WithDescription("error searching identity provider").WithCause(err, false)
	}
	if idp == nil {
		return nil, app_errors.NotFoundError(identity_providers.ErrIdentityProviderNotFound).
			WithDescription("identity provider not found")
	}
	if !idp.IsEnabled() {
		return nil, app_errors.ForbiddenError(identity_providers.ErrIdentityProviderDisabled).
			WithDescription("identity provider is disabled")
	}
	return idp, nil
}

// FindEnabled returns the provider only if it exists and is enabled.
func (f *IdentityProviderFinder) FindEnabled(ctx context.Context, id ids.IdentityProviderId) (*identity_providers.IdentityProvider, error) {
	idp, err := f.repository.FindById(ctx, id)
	if err != nil {
		return nil, app_errors.InternalError(identity_providers.ErrIdentityProviderNotFound).
			WithDescription("error searching identity provider").WithCause(

			err, false)

	}
	if idp == nil {
		return nil, app_errors.NotFoundError(identity_providers.ErrIdentityProviderNotFound).
			WithDescription("identity provider not found")
	}
	if !idp.IsEnabled() {
		return nil, app_errors.ForbiddenError(identity_providers.ErrIdentityProviderDisabled).
			WithDescription("identity provider is disabled")
	}
	return idp, nil
}
