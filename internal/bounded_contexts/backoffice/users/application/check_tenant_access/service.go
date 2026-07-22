package check_tenant_access

import (
	"context"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/users"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/users/domain_services"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
)

type Service struct {
	finder *domain_services.UserFinder
}

func NewService(finder *domain_services.UserFinder) *Service {
	return &Service{finder: finder}
}

func (s *Service) CheckTenantAccess(ctx context.Context, qry *Query) (any, error) {
	tenant, err := s.finder.Find(ctx, qry.tenantId)
	if err != nil {
		return nil, app_errors.AddTrace(err)
	}

	var allowed bool
	switch {
	case qry.author.Role.IsAdmin():
		allowed = tenant.IsOwnedByAdmin(qry.author)
	case qry.author.Role.IsTenant():
		allowed = tenant.IsSelf(qry.author)
	}

	if !allowed {
		return nil, app_errors.ForbiddenError(users.ErrInvalidOwnership).
			WithDescription("caller does not have access to this tenant")
	}

	return users.NewOwnership(qry.tenantId), nil
}
