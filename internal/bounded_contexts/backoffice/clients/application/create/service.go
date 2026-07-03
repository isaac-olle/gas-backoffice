package create

import (
	"context"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients"
)

type Service struct {
	repository clients.Repository
	eventBus   event.Bus
}

func NewService(repository clients.Repository, eventBus event.Bus) *Service {
	return &Service{repository: repository, eventBus: eventBus}
}

func (s *Service) Create(ctx context.Context, cmd *Command) error {
	client, err := clients.NewClient(cmd.clientID, cmd.name, cmd.ownership, cmd.redirectUris, cmd.allowedFlows)
	if err != nil {
		return app_errors.BadRequestError(clients.ErrClientNotCreated).
			WithDescription("failed to create client").
			WithCause(err)
	}

	if err = s.repository.Save(ctx, client); err != nil {
		return app_errors.InternalError(clients.ErrClientNotSaved).
			WithDescription("failed to save client").
			WithCause(err)
	}

	s.eventBus.Publish(ctx, client.PullDomainEvents())
	return nil
}
