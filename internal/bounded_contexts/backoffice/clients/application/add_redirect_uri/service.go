package add_redirect_uri

import (
	"context"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients/vos"
)

type Service struct {
	repository clients.Repository
	eventBus   event.Bus
}

func NewService(repository clients.Repository, eventBus event.Bus) *Service {
	return &Service{repository: repository, eventBus: eventBus}
}

func (s *Service) Add(ctx context.Context, cmd *Command) error {
	uri, err := vos.NewRedirectUri(cmd.redirectUri)
	if err != nil {
		return err
	}

	client, err := s.repository.SearchById(ctx, cmd.clientId)
	if err != nil {
		return err
	}

	client.AddRedirectUri(uri)

	if err = s.repository.Save(ctx, client); err != nil {
		return app_errors.InternalError(clients.ErrClientNotSaved).
			WithDescription("failed to save client").
			WithCause(err)
	}

	s.eventBus.Publish(ctx, client.PullDomainEvents())
	return nil
}
