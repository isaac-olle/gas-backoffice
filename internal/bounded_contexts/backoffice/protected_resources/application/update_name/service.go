package update_name

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
	prFinder             *domain_services.ProtectedResourceFinder
	duplicateNameChecker *domain_services.DuplicateNameChecker
	authorizer           *sharedDomainServices.SubscriptionAuthorizer
	repository           protected_resources.Repository
	eventBus             event.Bus
}

func NewService(
	prFinder *domain_services.ProtectedResourceFinder,
	duplicateNameChecker *domain_services.DuplicateNameChecker,
	repository protected_resources.Repository,
	eventBus event.Bus,
	authorizer *sharedDomainServices.SubscriptionAuthorizer,
) *Service {
	return &Service{prFinder: prFinder, duplicateNameChecker: duplicateNameChecker, repository: repository, eventBus: eventBus, authorizer: authorizer}
}

func (s *Service) UpdateName(ctx context.Context, cmd *Command) error {
	newName, err := vos.NewProtectedResourceName(cmd.name)
	if err != nil {
		return app_errors.BadRequestError(protected_resources.ErrProtectedResourceNotUpdated).
			WithDescription("invalid protected resource name").WithCause(err, false)
	}

	pr, err := s.prFinder.Find(ctx, cmd.protectedResourceId)
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

	// Uniqueness is scoped to the resource's owner, which may be one of the
	// caller's tenants when the caller is their admin.
	if err = s.duplicateNameChecker.Check(ctx, cmd.name, pr.ToPrimitives().SubscriptionId); err != nil {
		return app_errors.AddTrace(err)
	}

	pr.UpdateName(newName)

	if err = s.repository.Save(ctx, pr); err != nil {
		if errors.Is(err, persistance.ErrConcurrentModification) {
			return app_errors.ConflictError(protected_resources.ErrProtectedResourceConcurrentModification).
				WithDescription("protected resource was modified by another request; reload and retry").
				WithCause(err, false)
		}
		return app_errors.InternalError(protected_resources.ErrProtectedResourceNotSaved).
			WithDescription("failed to save protected resource").WithCause(err, false)
	}

	s.eventBus.Publish(ctx, pr.PullDomainEvents())
	return nil
}
