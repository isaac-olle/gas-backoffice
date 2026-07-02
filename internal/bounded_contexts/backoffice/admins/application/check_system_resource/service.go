package check_system_resource

import (
	"context"
	"github.com/isaac-olle/gas/internal/shared/app_errors"

	"github.com/isaac-olle/gas/internal/bounded_contexts/backoffice/admins"
	"github.com/isaac-olle/gas/internal/bounded_contexts/backoffice/admins/domain_services"
)

type Service struct {
	finder *domain_services.AdminFinder
}

func NewService(finder *domain_services.AdminFinder) *Service {
	return &Service{finder: finder}
}

func (s *Service) Check(ctx context.Context, qry *Query) (any, error) {
	admin, err := s.finder.Find(ctx, qry.adminId)
	if err != nil {
		return nil, err
	}

	if admin.IsSystemResource() {
		return nil, app_errors.ForbiddenError(admins.ErrSystemResource).
			WithDescription("admin is a system resource and cannot be modified")
	}
	return false, nil
}
