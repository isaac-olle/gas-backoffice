package update_name

import (
	"context"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/protected_resources/domain_services"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/protected_resources/vos"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/protected_resources"
)

type Service struct {
	prFinder   *domain_services.ProtectedResourceFinder
	repository protected_resources.Repository
	eventBus   event.Bus
}

func NewService(prFinder *domain_services.ProtectedResourceFinder, repository protected_resources.Repository, eventBus event.Bus) *Service {
	return &Service{prFinder: prFinder, repository: repository, eventBus: eventBus}
}

func (s *Service) UpdateName(ctx context.Context, cmd *Command) error {
	newName, err := vos.NewProtectedResourceName(cmd.name)
	if err != nil {
		return app_errors.BadRequestError(protected_resources.ErrProtectedResourceNotUpdated).
			WithDescription("invalid protected resource name").
			WithCause(err)
	}
	pr, err := s.prFinder.Find(ctx, cmd.protectedResourceId)
	if err != nil {
		return app_errors.NotFoundError(protected_resources.ErrProtectedResourceNotFound).
			WithDescription("protected resource not found").
			WithCause(err)
	}

	pr.UpdateName(newName)

	if err = s.repository.Save(ctx, pr); err != nil {
		return app_errors.InternalError(protected_resources.ErrProtectedResourceNotSaved).
			WithDescription("failed to save protected resource").
			WithCause(err)
	}

	s.eventBus.Publish(ctx, pr.PullDomainEvents())
	return nil
}
