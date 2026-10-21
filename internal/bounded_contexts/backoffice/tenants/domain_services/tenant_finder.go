package domain_services

import (
	"context"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/tenants"
)

type TenantFinder struct {
	repository tenants.Repository
}

func NewTenantFinder(repository tenants.Repository) *TenantFinder {
	return &TenantFinder{repository: repository}
}

func (f *TenantFinder) Find(ctx context.Context, id ids.UserId) (*tenants.Tenant, error) {
	tenant, err := f.repository.SearchById(ctx, id)
	if err != nil {
		return nil, app_errors.InternalError(tenants.ErrUserNotFound).
			WithDescription("error searching tenant").WithCause(

			err, false)

	}
	if tenant == nil {
		return nil, app_errors.NotFoundError(tenants.ErrUserNotFound).
			WithDescription("tenant not found")
	}
	return tenant, nil
}
