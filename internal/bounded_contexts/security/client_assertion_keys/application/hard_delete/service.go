package hard_delete

import (
	"context"
	"github.com/isaac-olle/gas/internal/shared/app_errors"
	"github.com/isaac-olle/gas/internal/shared/domain/vo/ids"

	"github.com/isaac-olle/gas/internal/bounded_contexts/security/client_assertion_keys"
)

type Service struct {
	repository client_assertion_keys.Repository
}

func NewService(repository client_assertion_keys.Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) HardDelete(ctx context.Context, cmd *Command) error {
	id, err := ids.NewClientAssertionKeyIdFromString(cmd.id)
	if err != nil {
		return app_errors.BadRequestError(client_assertion_keys.ErrClientAssertionKeyNotFound).
			WithDescription("invalid client assertion key id").
			WithCause(err)
	}

	if err = s.repository.HardDelete(ctx, id); err != nil {
		return app_errors.InternalError(client_assertion_keys.ErrClientAssertionKeyNotSaved).
			WithDescription("failed to hard delete client assertion key").
			WithCause(err)
	}
	return nil
}
