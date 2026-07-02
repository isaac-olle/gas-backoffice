package hard_delete

import (
	"context"
	"github.com/isaac-olle/gas/internal/shared/app_errors"
	"github.com/isaac-olle/gas/internal/shared/domain/vo/ids"

	"github.com/isaac-olle/gas/internal/bounded_contexts/security/client_secrets"
)

type Service struct {
	repository client_secrets.Repository
}

func NewService(repository client_secrets.Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) HardDelete(ctx context.Context, cmd *Command) error {
	id, err := ids.NewClientSecretIdFromString(cmd.id)
	if err != nil {
		return app_errors.BadRequestError(client_secrets.ErrClientSecretNotFound).
			WithDescription("invalid client secret id").
			WithCause(err)
	}

	if err = s.repository.HardDelete(ctx, id); err != nil {
		return app_errors.InternalError(client_secrets.ErrClientSecretNotSaved).
			WithDescription("failed to hard delete client secret").
			WithCause(err)
	}
	return nil
}
