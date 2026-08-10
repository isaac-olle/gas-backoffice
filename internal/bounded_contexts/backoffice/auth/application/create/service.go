package create

import (
	"context"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/auth"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"
)

type Service struct {
	repository auth.Repository
	eventBus   event.Bus
}

func NewService(repository auth.Repository, eventBus event.Bus) *Service {
	return &Service{repository: repository, eventBus: eventBus}
}

func (s *Service) Create(ctx context.Context, cmd *Command) error {
	session, err := auth.NewSession(cmd.id, cmd.userId)
	if err != nil {
		return app_errors.BadRequestError(auth.ErrSessionNotCreated).
			WithDescription("failed to create session").
			WithCause(err)
	}

	if err = s.repository.Save(ctx, session); err != nil {
		return app_errors.InternalError(auth.ErrSessionNotSaved).
			WithDescription("failed to save session").
			WithCause(err)
	}

	s.eventBus.Publish(ctx, session.PullDomainEvents())
	return nil
}
