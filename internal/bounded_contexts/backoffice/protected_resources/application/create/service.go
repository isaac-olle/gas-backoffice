package create

import (
	"context"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/protected_resources"
	prDomainServices "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/protected_resources/domain_services"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"
)

type Service struct {
	duplicateNameChecker *prDomainServices.DuplicateNameChecker
	repository           protected_resources.Repository
	eventBus             event.Bus
}

func NewService(duplicateNameChecker *prDomainServices.DuplicateNameChecker, repository protected_resources.Repository, eventBus event.Bus) *Service {
	return &Service{duplicateNameChecker: duplicateNameChecker, repository: repository, eventBus: eventBus}
}

func (s *Service) Create(ctx context.Context, cmd *Command) error {
	pr, err := protected_resources.NewProtectedResource(cmd.id, cmd.name, cmd.ownership, cmd.isSystem)
	if err != nil {
		return app_errors.BadRequestError(protected_resources.ErrProtectedResourceNotCreated).
			WithDescription("failed to create protected resource").WithCause(err, false)
	}

	if err = s.duplicateNameChecker.Check(ctx, cmd.name, cmd.ownership.SubscriptionId.String()); err != nil {
		return app_errors.AddTrace(err)
	}

	if err = s.repository.Save(ctx, pr); err != nil {
		return app_errors.InternalError(protected_resources.ErrProtectedResourceNotSaved).
			WithDescription("failed to save protected resource").WithCause(err, false)
	}

	s.eventBus.Publish(ctx, pr.PullDomainEvents())
	return nil
}
