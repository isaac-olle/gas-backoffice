package _event_handlers

import (
	"context"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/protected_resources"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/criteria"
	backofficeevents "github.com/isaac-olle/gas-backoffice/internal/shared/events/backoffice"
)

// DisableProtectedResourcesOnUserDeletedEventHandler disables (does not
// delete) every protected resource owned by a user once that user is
// deleted, so they stop being usable without losing their data.
type DisableProtectedResourcesOnUserDeletedEventHandler struct {
	repository protected_resources.Repository
	eventBus   event.Bus
}

func NewDisableProtectedResourcesOnUserDeletedEventHandler(repository protected_resources.Repository, eventBus event.Bus) *DisableProtectedResourcesOnUserDeletedEventHandler {
	return &DisableProtectedResourcesOnUserDeletedEventHandler{repository: repository, eventBus: eventBus}
}

func (h *DisableProtectedResourcesOnUserDeletedEventHandler) Handle(ctx context.Context, evt backofficeevents.UserDeleted) error {
	userId := evt.AggregateID()

	affectedPRs, err := h.repository.Matching(ctx, criteria.Criteria{
		Filters: []criteria.Filter{{Field: "user_id", Operator: criteria.Equals, Value: userId}},
	})
	if err != nil {
		return app_errors.InternalError(protected_resources.ErrProtectedResourceNotSaved).
			WithDescription("failed to search protected resources of the deleted user").WithCause(err, false)
	}

	for _, pr := range affectedPRs {
		if !pr.IsEnabled() {
			continue
		}
		pr.Disable()
		if err = h.repository.Save(ctx, pr); err != nil {
			return app_errors.InternalError(protected_resources.ErrProtectedResourceNotSaved).
				WithDescription("failed to disable protected resource of deleted user").WithCause(err, false)
		}
		h.eventBus.Publish(ctx, pr.PullDomainEvents())
	}

	return nil
}
