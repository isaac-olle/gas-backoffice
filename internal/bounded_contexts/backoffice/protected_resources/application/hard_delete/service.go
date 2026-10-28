package hard_delete

import (
	"context"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"
	sharedDomainServices "github.com/isaac-olle/gas-backoffice/internal/shared/domain_services"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/protected_resources"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/protected_resources/domain_services"
)

type Service struct {
	prFinder   *domain_services.ProtectedResourceFinder
	repository protected_resources.Repository
	eventBus   event.Bus
	authorizer *sharedDomainServices.SubscriptionAuthorizer
}

func NewService(prFinder *domain_services.ProtectedResourceFinder, repository protected_resources.Repository, eventBus event.Bus, authorizer *sharedDomainServices.SubscriptionAuthorizer) *Service {
	return &Service{prFinder: prFinder, repository: repository, eventBus: eventBus, authorizer: authorizer}
}

func (s *Service) HardDelete(ctx context.Context, cmd *Command) error {
	pr, err := s.prFinder.Find(ctx, cmd.protectedResourceId)
	if err != nil {
		return app_errors.AddTrace(err)
	}

	if err = s.authorizer.Check(pr, cmd.author); err != nil {
		return app_errors.AddTrace(err)
	}

	if pr.IsSystemResource() {
		return app_errors.ForbiddenError(protected_resources.ErrSystemResource).
			WithDescription("protected resource is a system resource and cannot be deleted")
	}

	pr.Delete()

	if err = s.repository.HardDelete(ctx, cmd.protectedResourceId); err != nil {
		return app_errors.InternalError(protected_resources.ErrProtectedResourceNotSaved).WithCause(
			err, false)

	}

	s.eventBus.Publish(ctx, pr.PullDomainEvents())
	return nil
}
