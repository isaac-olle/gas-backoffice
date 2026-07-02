package check_system_resource

import (
	"context"
	"github.com/isaac-olle/gas/internal/shared/app_errors"

	"github.com/isaac-olle/gas/internal/bounded_contexts/backoffice/protected_resources"
	"github.com/isaac-olle/gas/internal/bounded_contexts/backoffice/protected_resources/domain_services"
)

type Service struct {
	finder *domain_services.ProtectedResourceFinder
}

func NewService(finder *domain_services.ProtectedResourceFinder) *Service {
	return &Service{finder: finder}
}

func (s *Service) Check(ctx context.Context, qry *Query) (any, error) {
	pr, err := s.finder.Find(ctx, qry.protectedResourceId)
	if err != nil {
		return nil, err
	}

	if pr.IsSystemResource() {
		return nil, app_errors.ForbiddenError(protected_resources.ErrSystemResource).
			WithDescription("protected resource is a system resource and cannot be modified")
	}
	return false, nil
}
