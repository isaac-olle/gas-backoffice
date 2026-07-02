package check_ownership

import (
	"context"
	"github.com/isaac-olle/gas/internal/shared/app_errors"
	"github.com/isaac-olle/gas/internal/shared/domain/vo/ids"

	"github.com/isaac-olle/gas/internal/bounded_contexts/backoffice/tenants"
	"github.com/isaac-olle/gas/internal/bounded_contexts/backoffice/tenants/domain_services"
)

type Service struct {
	finder *domain_services.TenantFinder
}

func NewService(finder *domain_services.TenantFinder) *Service {
	return &Service{finder: finder}
}

func (s *Service) CheckOwnership(ctx context.Context, qry *Query) (any, error) {
	tenantId, err := ids.NewTenantIdFromString(qry.tenantId)
	if err != nil {
		return nil, app_errors.BadRequestError(tenants.ErrTenantNotFound).
			WithDescription("invalid tenant id").
			WithCause(err)
	}

	ownership := tenants.NewOwnership(qry.adminId)

	tenant, err := s.finder.Find(ctx, tenantId)
	if err != nil {
		return nil, app_errors.AddTrace(err)
	}

	if !tenant.CheckOwnership(tenants.NewOwnership(qry.adminId)) {
		return nil, app_errors.ForbiddenError(tenants.ErrInvalidOwnership).
			WithDescription("tenant does not belong to the provided admin")
	}

	return ownership, nil
}
