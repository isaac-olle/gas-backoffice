package get

import (
	"context"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/users/domain_services"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
)

type Service struct {
	userFinder *domain_services.UserFinder
}

func NewService(userFinder *domain_services.UserFinder) *Service {
	return &Service{userFinder: userFinder}
}

func (s *Service) Get(ctx context.Context, qry *Query) (any, error) {
	user, err := s.userFinder.Find(ctx, qry.userId)
	if err != nil {
		return nil, app_errors.AddTrace(err)
	}

	if err = user.EnsureUserAccess(qry.author); err != nil {
		return nil, app_errors.AddTrace(err)
	}

	return user, nil
}
