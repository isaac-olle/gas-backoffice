package update_name

import (
	"context"
	"errors"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/subscriptions"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"
	"github.com/isaac-olle/gas-backoffice/internal/shared/persistance"
)

type Service struct {
	repository subscriptions.Repository
	eventBus   event.Bus
}

func NewService(repository subscriptions.Repository, eventBus event.Bus) *Service {
	return &Service{repository: repository, eventBus: eventBus}
}

func (s *Service) UpdateName(ctx context.Context, cmd *Command) error {
	subscription, err := s.repository.SearchById(ctx, cmd.subscriptionId)
	if err != nil {
		return app_errors.InternalError(subscriptions.ErrSubscriptionNotFound).
			WithDescription("failed to find subscription").WithCause(err, false)
	}
	if subscription == nil {
		// PUT is an idempotent upsert: when the client supplies a new ID,
		// create the subscription owned by the authenticated author.
		subscription, err = subscriptions.NewSubscription(cmd.subscriptionId, cmd.name, cmd.ownerId, false)
		if err != nil {
			return app_errors.AddTrace(err)
		}
		if err = s.repository.Save(ctx, subscription); err != nil {
			if errors.Is(err, persistance.ErrConcurrentModification) {
				return app_errors.ConflictError(subscriptions.ErrSubscriptionConcurrentModification).
					WithDescription("subscription was modified by another request; reload and retry").WithCause(err, false)
			}
			return app_errors.InternalError(subscriptions.ErrSubscriptionNotSaved).
				WithDescription("failed to save subscription").WithCause(err, false)
		}
		s.eventBus.Publish(ctx, subscription.PullDomainEvents())
		return nil
	}
	if subscription.IsSystemResource() {
		return app_errors.ForbiddenError(subscriptions.ErrSystemResource).
			WithDescription("subscription is a system resource and cannot be modified")
	}
	if err = subscription.ChangeName(cmd.name); err != nil {
		return app_errors.AddTrace(err)
	}
	if err = s.repository.Save(ctx, subscription); err != nil {
		if errors.Is(err, persistance.ErrConcurrentModification) {
			return app_errors.ConflictError(subscriptions.ErrSubscriptionConcurrentModification).
				WithDescription("subscription was modified by another request; reload and retry").WithCause(err, false)
		}
		return app_errors.InternalError(subscriptions.ErrSubscriptionNotSaved).
			WithDescription("failed to save subscription").WithCause(err, false)
	}
	s.eventBus.Publish(ctx, subscription.PullDomainEvents())
	return nil
}
