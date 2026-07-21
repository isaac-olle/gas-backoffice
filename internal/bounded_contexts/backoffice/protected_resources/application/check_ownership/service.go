package check_ownership

import (
	"context"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/protected_resources"
)

type Service struct {
	repository protected_resources.Repository
}

func NewService(repository protected_resources.Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) CheckOwnership(ctx context.Context, qry *Query) (any, error) {
	userId, err := ids2.NewUserIdFromString(qry.userId)
	if err != nil {
		return nil, app_errors.BadRequestError(protected_resources.ErrInvalidUserId).
			WithDescription("invalid user id").
			WithCause(err)
	}

	ownership := protected_resources.NewOwnership(userId)

	prId, err := ids2.NewProtectedResourceIdFromString(qry.protectedResourceId)
	if err != nil {
		return nil, app_errors.BadRequestError(protected_resources.ErrProtectedResourceNotFound).
			WithDescription("invalid protected resource id").
			WithCause(err)
	}

	pr, err := s.repository.SearchById(ctx, prId)
	if err != nil {
		return nil, app_errors.NotFoundError(protected_resources.ErrProtectedResourceNotFound).
			WithDescription("protected resource not found").
			WithCause(err)
	}

	if !pr.CheckOwnership(ownership) {
		return nil, app_errors.ForbiddenError(protected_resources.ErrInvalidOwnership).
			WithDescription("protected resource does not belong to the provided user")
	}

	return ownership, nil
}
