package get

import (
	"context"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/tenants/domain_services"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
)

type Service struct {
	tenantFinder *domain_services.TenantFinder
}

func NewService(tenantFinder *domain_services.TenantFinder) *Service {
	return &Service{tenantFinder: tenantFinder}
}

func (s *Service) Get(ctx context.Context, qry *Query) (any, error) {
	tenant, err := s.tenantFinder.Find(ctx, qry.userId)
	if err != nil {
		return nil, app_errors.AddTrace(err)
	}

	if err = tenant.EnsureTenantAccess(qry.author); err != nil {
		return nil, app_errors.AddTrace(err)
	}

	return tenant, nil
}
