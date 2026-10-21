package _event_handlers

import (
	"context"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/scopes"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	backofficeevents "github.com/isaac-olle/gas-backoffice/internal/shared/events/backoffice"
)

// RemoveClientFromScopesOnClientDeletedEventHandler keeps Scope.allowedClientIds
// free of dangling references: when a client is deleted, it's removed from
// every scope that still lists it as an allowed client.
type RemoveClientFromScopesOnClientDeletedEventHandler struct {
	repository scopes.Repository
	eventBus   event.Bus
}

func NewRemoveClientFromScopesOnClientDeletedEventHandler(repository scopes.Repository, eventBus event.Bus) *RemoveClientFromScopesOnClientDeletedEventHandler {
	return &RemoveClientFromScopesOnClientDeletedEventHandler{repository: repository, eventBus: eventBus}
}

func (h *RemoveClientFromScopesOnClientDeletedEventHandler) Handle(ctx context.Context, evt backofficeevents.ClientDeleted) error {
	clientId, err := ids2.NewClientIdFromString(evt.AggregateID())
	if err != nil {
		return app_errors.BadRequestError(scopes.ErrScopeNotFound).
			WithDescription("invalid client id in ClientDeleted event").WithCause(err, false)
	}

	affectedScopes, err := h.repository.FindByAllowedClientId(ctx, clientId)
	if err != nil {
		return app_errors.InternalError(scopes.ErrScopeNotSaved).
			WithDescription("failed to search scopes allowing the deleted client").WithCause(err, false)
	}

	for _, scope := range affectedScopes {
		scope.RemoveAllowedClient(clientId)
		if err = h.repository.Save(ctx, scope); err != nil {
			return app_errors.InternalError(scopes.ErrScopeNotSaved).
				WithDescription("failed to remove deleted client from scope").WithCause(err, false)
		}
		h.eventBus.Publish(ctx, scope.PullDomainEvents())
	}

	return nil
}
