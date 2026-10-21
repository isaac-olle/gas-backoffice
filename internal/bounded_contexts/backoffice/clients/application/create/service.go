package create

import (
	"context"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients"
	clientDomainServices "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients/domain_services"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"
)

type Service struct {
	duplicateNameChecker *clientDomainServices.DuplicateNameChecker
	repository           clients.Repository
	eventBus             event.Bus
}

func NewService(duplicateNameChecker *clientDomainServices.DuplicateNameChecker, repository clients.Repository, eventBus event.Bus) *Service {
	return &Service{duplicateNameChecker: duplicateNameChecker, repository: repository, eventBus: eventBus}
}

func (s *Service) Create(ctx context.Context, cmd *Command) error {
	client, err := clients.NewClient(cmd.clientID, cmd.name, cmd.ownership, cmd.redirectUris, cmd.allowedGrantTypes, cmd.public, cmd.isSystem)
	if err != nil {
		return app_errors.BadRequestError(clients.ErrClientNotCreated).
			WithDescription("failed to create client").WithCause(err, false)
	}

	if err = s.duplicateNameChecker.Check(ctx, cmd.name, cmd.ownership.SubscriptionId.String()); err != nil {
		return app_errors.AddTrace(err)
	}

	if err = s.repository.Save(ctx, client); err != nil {
		return app_errors.InternalError(clients.ErrClientNotSaved).
			WithDescription("failed to save client").WithCause(err, false)
	}

	s.eventBus.Publish(ctx, client.PullDomainEvents())
	return nil
}
