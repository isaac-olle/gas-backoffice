package create

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

func NewService(r subscriptions.Repository, b event.Bus) *Service {
	return &Service{repository: r, eventBus: b}
}
func (s *Service) Create(ctx context.Context, cmd *Command) error {
	aggregate, err := subscriptions.NewSubscription(cmd.id, cmd.name, cmd.ownerId, cmd.isSystem)
	if err != nil {
		return app_errors.AddTrace(err)
	}
	if err = s.repository.Save(ctx, aggregate); err != nil {
		if errors.Is(err, persistance.ErrConcurrentModification) {
			return app_errors.ConflictError(subscriptions.ErrSubscriptionConcurrentModification).WithCause(err, false)
		}
		return app_errors.InternalError(subscriptions.ErrSubscriptionNotSaved).WithCause(err, false)
	}
	s.eventBus.Publish(ctx, aggregate.PullDomainEvents())
	return nil
}
