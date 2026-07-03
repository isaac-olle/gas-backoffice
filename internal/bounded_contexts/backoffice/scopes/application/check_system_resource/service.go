package check_system_resource

import (
	"context"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/scopes"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/scopes/domain_services"
)

type Service struct {
	finder *domain_services.ScopeFinder
}

func NewService(finder *domain_services.ScopeFinder) *Service {
	return &Service{finder: finder}
}

func (s *Service) Check(ctx context.Context, qry *Query) (any, error) {
	scope, err := s.finder.Find(ctx, qry.scopeId)
	if err != nil {
		return nil, err
	}

	if scope.IsSystemResource() {
		return nil, app_errors.ForbiddenError(scopes.ErrSystemResource).
			WithDescription("scope is a system resource and cannot be modified")
	}
	return false, nil
}
