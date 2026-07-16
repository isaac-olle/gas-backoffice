package remove_allowed_flow

import (
	"context"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients/domain_services"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients/vos"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"
)

type Service struct {
	clientFinder *domain_services.ClientFinder
	repository   clients.Repository
	eventBus     event.Bus
}

func NewService(clientFinder *domain_services.ClientFinder, repository clients.Repository, eventBus event.Bus) *Service {
	return &Service{clientFinder: clientFinder, repository: repository, eventBus: eventBus}
}

func (s *Service) Remove(ctx context.Context, cmd *Command) error {
	flow, err := vos.NewFlow(cmd.flow)
	if err != nil {
		return err
	}

	client, err := s.clientFinder.Find(ctx, cmd.clientId)
	if err != nil {
		return err
	}

	client.RemoveAllowedFlow(flow)

	if err = s.repository.Save(ctx, client); err != nil {
		return app_errors.InternalError(clients.ErrClientNotSaved).
			WithDescription("failed to save client").
			WithCause(err)
	}

	s.eventBus.Publish(ctx, client.PullDomainEvents())
	return nil
}
