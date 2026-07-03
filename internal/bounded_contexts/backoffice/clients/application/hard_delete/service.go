package hard_delete

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

func (s *Service) HardDelete(ctx context.Context, cmd *Command) error {
	clientId, err := ids.NewClientIdFromString(cmd.clientId)
	if err != nil {
		return app_errors.BadRequestError(clients.ErrClientNotFound).
			WithDescription("invalid client id").
			WithCause(err)
	}

	if err = s.repository.HardDelete(ctx, clientId); err != nil {
		return app_errors.InternalError(clients.ErrClientNotSaved).
			WithDescription("failed to hard delete client").
			WithCause(err)
	}
	return nil
}
