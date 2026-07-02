package add_allowed_flow

import (
	"context"
	"github.com/isaac-olle/gas/internal/shared/app_errors"
	"github.com/isaac-olle/gas/internal/shared/buses/event"

	"github.com/isaac-olle/gas/internal/bounded_contexts/backoffice/clients"
	"github.com/isaac-olle/gas/internal/bounded_contexts/backoffice/clients/vos"
)

type Service struct {
	repository clients.Repository
	eventBus   event.Bus
}

func NewService(repository clients.Repository, eventBus event.Bus) *Service {
	return &Service{repository: repository, eventBus: eventBus}
}

func (s *Service) Add(ctx context.Context, cmd *Command) error {
	flow, err := vos.NewFlow(cmd.flow)
	if err != nil {
		return err
	}

	client, err := s.repository.SearchById(ctx, cmd.clientId)
	if err != nil {
		return err
	}

	client.AddAllowedFlow(flow)

	if err = s.repository.Save(ctx, client); err != nil {
		return app_errors.InternalError(clients.ErrClientNotSaved).
			WithDescription("failed to save client").
			WithCause(err)
	}

	s.eventBus.Publish(ctx, client.PullDomainEvents())
	return nil
}
