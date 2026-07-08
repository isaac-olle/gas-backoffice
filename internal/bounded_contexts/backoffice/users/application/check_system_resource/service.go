package check_system_resource

import (
	"context"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/users"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/users/domain_services"
)

type Service struct {
	finder *domain_services.UserFinder
}

func NewService(finder *domain_services.UserFinder) *Service {
	return &Service{finder: finder}
}

func (s *Service) Check(ctx context.Context, qry *Query) (any, error) {
	user, err := s.finder.Find(ctx, qry.userId)
	if err != nil {
		return nil, err
	}

	if user.IsSystemResource() {
		return nil, app_errors.ForbiddenError(users.ErrSystemResource).
			WithDescription("user is a system resource and cannot be modified")
	}
	return false, nil
}
