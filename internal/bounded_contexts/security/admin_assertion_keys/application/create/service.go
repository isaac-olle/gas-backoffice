package create

import (
	"context"
	"github.com/isaac-olle/gas/internal/shared/app_errors"
	"github.com/isaac-olle/gas/internal/shared/buses/event"
	"github.com/isaac-olle/gas/internal/shared/domain/vo/ids"

	"github.com/isaac-olle/gas/internal/bounded_contexts/security/admin_assertion_keys"
)

type Service struct {
	repository admin_assertion_keys.Repository
	eventBus   event.Bus
}

func NewService(repository admin_assertion_keys.Repository, eventBus event.Bus) *Service {
	return &Service{repository: repository, eventBus: eventBus}
}

func (s *Service) Create(ctx context.Context, cmd *Command) error {
	adminId, err := ids.NewAdminIdFromString(cmd.adminId)
	if err != nil {
		return app_errors.BadRequestError(admin_assertion_keys.ErrAdminAssertionKeyNotCreated).
			WithDescription("invalid admin id").
			WithCause(err)
	}

	aak, err := admin_assertion_keys.NewAdminAssertionKey(cmd.id, adminId, cmd.algorithm, cmd.publicKey)
	if err != nil {
		return app_errors.BadRequestError(admin_assertion_keys.ErrAdminAssertionKeyNotCreated).
			WithDescription("failed to create admin assertion key").
			WithCause(err)
	}

	if err = s.repository.Save(ctx, aak); err != nil {
		return app_errors.InternalError(admin_assertion_keys.ErrAdminAssertionKeyNotSaved).
			WithDescription("failed to save admin assertion key").
			WithCause(err)
	}

	s.eventBus.Publish(ctx, aak.PullDomainEvents())
	return nil
}
