package check_ownership

import (
	"context"
	"github.com/isaac-olle/gas/internal/shared/app_errors"
	ids2 "github.com/isaac-olle/gas/internal/shared/domain/vo/ids"

	"github.com/isaac-olle/gas/internal/bounded_contexts/backoffice/protected_resources"
)

type Service struct {
	repository protected_resources.Repository
}

func NewService(repository protected_resources.Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) CheckOwnership(ctx context.Context, qry *Query) (any, error) {
	tenantId, err := ids2.NewTenantIdFromString(qry.tenantId)
	if err != nil {
		return nil, app_errors.BadRequestError(protected_resources.ErrInvalidTenantId).
			WithDescription("invalid tenant id").
			WithCause(err)
	}

	ownership := protected_resources.NewOwnership(tenantId, qry.adminId)

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
			WithDescription("protected resource does not belong to the provided tenant")
	}

	return ownership, nil
}
