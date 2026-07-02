package check_system_resource

import (
	"context"
	"github.com/isaac-olle/gas/internal/shared/app_errors"

	"github.com/isaac-olle/gas/internal/bounded_contexts/backoffice/tenants"
	"github.com/isaac-olle/gas/internal/bounded_contexts/backoffice/tenants/domain_services"
)

type Service struct {
	finder *domain_services.TenantFinder
}

func NewService(finder *domain_services.TenantFinder) *Service {
	return &Service{finder: finder}
}

func (s *Service) Check(ctx context.Context, qry *Query) (any, error) {
	tenant, err := s.finder.Find(ctx, qry.tenantId)
	if err != nil {
		return nil, err
	}

	if tenant.IsSystemResource() {
		return nil, app_errors.ForbiddenError(tenants.ErrSystemResource).
			WithDescription("tenant is a system resource and cannot be modified")
	}
	return false, nil
}
