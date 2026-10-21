package _event_handlers

import (
	"context"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/criteria"
	backofficeevents "github.com/isaac-olle/gas-backoffice/internal/shared/events/backoffice"
)

// DisableClientsOnUserDeletedEventHandler disables (does not delete) every
// client owned by a user once that user is deleted, so they stop being
// usable without losing their data.
type DisableClientsOnUserDeletedEventHandler struct {
	repository clients.Repository
	eventBus   event.Bus
}

func NewDisableClientsOnUserDeletedEventHandler(repository clients.Repository, eventBus event.Bus) *DisableClientsOnUserDeletedEventHandler {
	return &DisableClientsOnUserDeletedEventHandler{repository: repository, eventBus: eventBus}
}

func (h *DisableClientsOnUserDeletedEventHandler) Handle(ctx context.Context, evt backofficeevents.UserDeleted) error {
	userId := evt.AggregateID()

	affectedClients, err := h.repository.Matching(ctx, criteria.Criteria{
		Filters: []criteria.Filter{{Field: "user_id", Operator: criteria.Equals, Value: userId}},
	})
	if err != nil {
		return app_errors.InternalError(clients.ErrClientNotSaved).
			WithDescription("failed to search clients of the deleted user").WithCause(err, false)
	}

	for _, client := range affectedClients {
		if !client.IsEnabled() {
			continue
		}
		client.Disable()
		if err = h.repository.Save(ctx, client); err != nil {
			return app_errors.InternalError(clients.ErrClientNotSaved).
				WithDescription("failed to disable client of deleted user").WithCause(err, false)
		}
		h.eventBus.Publish(ctx, client.PullDomainEvents())
	}

	return nil
}
