package _event_handlers

import (
	"context"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	
)

type AddOpenIDScopeOnClientCreatedEventHandler struct {
	repository    clients.Repository
	eventBus      event.Bus
	openIDScopeId func() (ids2.ScopeId, error)
}

func NewAddOpenIDScopeOnClientCreatedEventHandler(
	repository clients.Repository,
	eventBus event.Bus,
	openIDScopeId func() (ids2.ScopeId, error),
) *AddOpenIDScopeOnClientCreatedEventHandler {
	return &AddOpenIDScopeOnClientCreatedEventHandler{
		repository:    repository,
		eventBus:      eventBus,
		openIDScopeId: openIDScopeId,
	}
}

func (h *AddOpenIDScopeOnClientCreatedEventHandler) Handle(ctx context.Context, evt backofficeevents.ClientCreated) error {
	clientId, err := ids2.NewClientIdFromString(evt.AggregateID())
	if err != nil {
		return app_errors.BadRequestError(clients.ErrClientNotFound).
			WithDescription("invalid client id in ClientCreated event").
			WithCause(err)
	}

	client, err := h.repository.SearchById(ctx, clientId)
	if err != nil {
		return app_errors.InternalError(clients.ErrClientNotSaved).
			WithDescription("failed to find client when handling ClientCreated event").
			WithCause(err)
	}

	scopeId, err := h.openIDScopeId()
	if err != nil {
		return app_errors.InternalError("OPENID_SCOPE_ID_INVALID").
			WithDescription("invalid openid scope id").
			WithCause(err)
	}

	client.AddAllowedScope(scopeId)

	if err = h.repository.Save(ctx, client); err != nil {
		return app_errors.InternalError(clients.ErrClientNotSaved).
			WithDescription("failed to save client after adding openid scope").
			WithCause(err)
	}

	h.eventBus.Publish(ctx, client.PullDomainEvents())
	return nil
}
