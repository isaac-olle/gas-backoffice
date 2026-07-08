package check_ownership

import (
	"context"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients"
)

type Service struct {
	repository clients.Repository
}

func NewService(repository clients.Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) CheckOwnership(ctx context.Context, qry *Query) (any, error) {
	clientId, err := ids.NewClientIdFromString(qry.clientId)
	if err != nil {
		return nil, app_errors.BadRequestError(clients.ErrClientNotFound).
			WithDescription("invalid client id").
			WithCause(err)
	}

	ownership := clients.NewOwnership(qry.adminId, qry.userId)

	client, err := s.repository.SearchById(ctx, clientId)
	if err != nil {
		return nil, app_errors.NotFoundError(clients.ErrClientNotFound).
			WithDescription("client not found").
			WithCause(err)
	}

	if !client.CheckOwnership(ownership) {
		return nil, app_errors.ForbiddenError(clients.ErrInvalidOwnership).
			WithDescription("client does not belong to the provided admin and user")
	}

	return ownership.ToPrimitives(), nil
}
