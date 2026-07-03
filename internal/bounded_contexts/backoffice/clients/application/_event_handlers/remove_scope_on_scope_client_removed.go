package _event_handlers

import (
	"context"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	backofficeevents "github.com/isaac-olle/gas-backoffice/internal/shared/events/backoffice"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients"
)

type RemoveScopeOnScopeClientRemovedEventHandler struct {
	repository clients.Repository
	eventBus   event.Bus
}

func NewRemoveScopeOnScopeClientRemovedEventHandler(repository clients.Repository, eventBus event.Bus) *RemoveScopeOnScopeClientRemovedEventHandler {
	return &RemoveScopeOnScopeClientRemovedEventHandler{repository: repository, eventBus: eventBus}
}

func (h *RemoveScopeOnScopeClientRemovedEventHandler) Handle(ctx context.Context, evt backofficeevents.ScopeClientRemoved) error {
	clientId, err := ids2.NewClientIdFromString(evt.ClientId)
	if err != nil {
		return app_errors.BadRequestError(clients.ErrClientNotFound).
			WithDescription("invalid client id in ScopeClientRemoved event").
			WithCause(err)
	}

	client, err := h.repository.SearchById(ctx, clientId)
	if err != nil {
		return app_errors.NotFoundError(clients.ErrClientNotFound).
			WithDescription("client not found when handling ScopeClientRemoved event").
			WithCause(err)
	}

	scopeId, err := ids2.NewScopeIdFromString(evt.AggregateID())
	if err != nil {
		return err
	}

	client.RemoveAllowedScope(scopeId)

	if err = h.repository.Save(ctx, client); err != nil {
		return app_errors.InternalError(clients.ErrClientNotSaved).
			WithDescription("failed to save client after removing scope").
			WithCause(err)
	}

	h.eventBus.Publish(ctx, client.PullDomainEvents())
	return nil
}
