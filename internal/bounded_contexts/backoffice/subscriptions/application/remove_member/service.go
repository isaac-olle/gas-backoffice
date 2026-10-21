package remove_member

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
func (s *Service) Remove(ctx context.Context, cmd *Command) error {
	a, err := s.repository.SearchById(ctx, cmd.id)
	if err != nil {
		return app_errors.InternalError(subscriptions.ErrSubscriptionNotFound).WithCause(err, false)
	}
	if a == nil {
		return app_errors.NotFoundError(subscriptions.ErrSubscriptionNotFound)
	}
	a.RemoveMember(cmd.userId)
	if err = s.repository.Save(ctx, a); err != nil {
		if errors.Is(err, persistance.ErrConcurrentModification) {
			return app_errors.ConflictError(subscriptions.ErrSubscriptionConcurrentModification).WithCause(err, false)
		}
		return app_errors.InternalError(subscriptions.ErrSubscriptionNotSaved).WithCause(err, false)
	}
	s.eventBus.Publish(ctx, a.PullDomainEvents())
	return nil
}
