package domain_services

import (
	"context"
	"github.com/isaac-olle/gas/internal/shared/app_errors"
	"github.com/isaac-olle/gas/internal/shared/domain/vo/ids"

	"github.com/isaac-olle/gas/internal/bounded_contexts/backoffice/tenants"
)

type TenantFinder struct {
	repository tenants.Repository
}

func NewTenantFinder(repository tenants.Repository) *TenantFinder {
	return &TenantFinder{repository: repository}
}

func (f *TenantFinder) Find(ctx context.Context, id ids.TenantId) (*tenants.Tenant, error) {
	tenant, err := f.repository.SearchById(ctx, id)
	if err != nil {
		return nil, app_errors.InternalError(tenants.ErrTenantNotFound).
			WithDescription("error searching tenant").
			WithCause(err)
	}
	if tenant == nil {
		return nil, app_errors.NotFoundError(tenants.ErrTenantNotFound).
			WithDescription("tenant not found")
	}
	return tenant, nil
}
