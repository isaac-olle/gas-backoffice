package _event_handlers

import (
	"context"
	"errors"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients/domain_services"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	backofficeevents "github.com/isaac-olle/gas-backoffice/internal/shared/events/backoffice"
	"github.com/isaac-olle/gas-backoffice/internal/shared/persistance"
)

type RemoveAllowedScopeOnScopeClientRemovedEventHandler struct {
	clientFinder *domain_services.ClientFinder
	repository   clients.Repository
}

func NewRemoveAllowedScopeOnScopeClientRemovedEventHandler(
	clientFinder *domain_services.ClientFinder,
	repository clients.Repository,
) *RemoveAllowedScopeOnScopeClientRemovedEventHandler {
	return &RemoveAllowedScopeOnScopeClientRemovedEventHandler{clientFinder: clientFinder, repository: repository}
}

func (h *RemoveAllowedScopeOnScopeClientRemovedEventHandler) Handle(ctx context.Context, evt backofficeevents.ScopeClientRemoved) error {
	clientId, err := ids2.NewClientIdFromString(evt.ClientId)
	if err != nil {
		return app_errors.BadRequestError(clients.ErrClientNotFound).
			WithDescription("invalid client id in ScopeClientRemoved event").WithCause(err, false)
	}

	scopeId, err := ids2.NewScopeIdFromString(evt.AggregateID())
	if err != nil {
		return app_errors.BadRequestError(clients.ErrClientNotFound).
			WithDescription("invalid scope id in ScopeClientRemoved event").WithCause(err, false)
	}

	client, err := h.clientFinder.Find(ctx, clientId)
	if err != nil {
		return app_errors.AddTrace(err)
	}

	client.RemoveAllowedScope(scopeId)

	if err = h.repository.Save(ctx, client); err != nil {
		if errors.Is(err, persistance.ErrConcurrentModification) {
			return app_errors.ConflictError(clients.ErrClientConcurrentModification).
				WithDescription("client was modified concurrently; event will be logged as failed").WithCause(err, false)
		}
		return app_errors.InternalError(clients.ErrClientNotSaved).
			WithDescription("failed to save client").WithCause(err, false)
	}

	return nil
}
