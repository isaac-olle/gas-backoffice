package check_ownership

import (
	"context"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/users"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/users/domain_services"
)

type Service struct {
	finder *domain_services.UserFinder
}

func NewService(finder *domain_services.UserFinder) *Service {
	return &Service{finder: finder}
}

func (s *Service) CheckOwnership(ctx context.Context, qry *Query) (any, error) {
	userId, err := ids.NewUserIdFromString(qry.userId)
	if err != nil {
		return nil, app_errors.BadRequestError(users.ErrUserNotFound).
			WithDescription("invalid user id").
			WithCause(err)
	}

	user, err := s.finder.Find(ctx, userId)
	if err != nil {
		return nil, app_errors.AddTrace(err)
	}

	ownership := users.NewOwnership(userId)

	if !user.CheckOwnership(ownership) {
		return nil, app_errors.ForbiddenError(users.ErrInvalidOwnership).
			WithDescription("user not found with the provided id")
	}

	return ownership, nil
}
