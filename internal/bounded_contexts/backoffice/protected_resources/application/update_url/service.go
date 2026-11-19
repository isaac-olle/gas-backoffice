package update_url

import (
	"context"
	"errors"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/protected_resources"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/protected_resources/domain_services"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/protected_resources/vos"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"
	sharedDomainServices "github.com/isaac-olle/gas-backoffice/internal/shared/domain_services"
	"github.com/isaac-olle/gas-backoffice/internal/shared/persistance"
)

type Service struct {
	finder     *domain_services.ProtectedResourceFinder
	repository protected_resources.Repository
	eventBus   event.Bus
	authorizer *sharedDomainServices.SubscriptionAuthorizer
}

func NewService(
	finder *domain_services.ProtectedResourceFinder,
	repository protected_resources.Repository,
	eventBus event.Bus,
	authorizer *sharedDomainServices.SubscriptionAuthorizer,
) *Service {
	return &Service{finder: finder, repository: repository, eventBus: eventBus, authorizer: authorizer}
}

func (s *Service) UpdateURL(ctx context.Context, cmd *Command) error {
	newURL, err := vos.NewProtectedResourceURL(cmd.url)
	if err != nil {
		return app_errors.BadRequestError(protected_resources.ErrProtectedResourceNotUpdated).
			WithDescription("invalid protected resource URL").WithCause(err, false)
	}
	pr, err := s.finder.Find(ctx, cmd.protectedResourceId)
	if err != nil {
		return app_errors.AddTrace(err)
	}
	if err = s.authorizer.Check(pr, cmd.author); err != nil {
		return app_errors.AddTrace(err)
	}
	if pr.IsSystemResource() {
		return app_errors.ForbiddenError(protected_resources.ErrSystemResource).
			WithDescription("protected resource is a system resource and cannot be modified")
	}
	pr.UpdateURL(newURL)
	if err = s.repository.Save(ctx, pr); err != nil {
		if errors.Is(err, persistance.ErrConcurrentModification) {
			return app_errors.ConflictError(protected_resources.ErrProtectedResourceConcurrentModification).
				WithDescription("protected resource was modified by another request; reload and retry").WithCause(err, false)
		}
		return app_errors.InternalError(protected_resources.ErrProtectedResourceNotSaved).
			WithDescription("failed to save protected resource").WithCause(err, false)
	}
	s.eventBus.Publish(ctx, pr.PullDomainEvents())
	return nil
}
