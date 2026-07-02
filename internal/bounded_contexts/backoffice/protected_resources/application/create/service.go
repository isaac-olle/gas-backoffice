package create

import (
	"context"
	"github.com/isaac-olle/gas/internal/shared/app_errors"
	"github.com/isaac-olle/gas/internal/shared/buses/event"

	"github.com/isaac-olle/gas/internal/bounded_contexts/backoffice/protected_resources"
)

type Service struct {
	repository protected_resources.Repository
	eventBus   event.Bus
}

func NewService(repository protected_resources.Repository, eventBus event.Bus) *Service {
	return &Service{repository: repository, eventBus: eventBus}
}

func (s *Service) Create(ctx context.Context, cmd *Command) error {
	pr, err := protected_resources.NewProtectedResource(cmd.id, cmd.name, cmd.ownership)
	if err != nil {
		return app_errors.BadRequestError(protected_resources.ErrProtectedResourceNotCreated).
			WithDescription("failed to create protected resource").
			WithCause(err)
	}

	if err = s.repository.Save(ctx, pr); err != nil {
		return app_errors.InternalError(protected_resources.ErrProtectedResourceNotSaved).
			WithDescription("failed to save protected resource").
			WithCause(err)
	}

	s.eventBus.Publish(ctx, pr.PullDomainEvents())
	return nil
}
