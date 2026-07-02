package _event_handlers

import (
	"context"
	"github.com/isaac-olle/gas/internal/shared/app_errors"
	"github.com/isaac-olle/gas/internal/shared/buses/event"
	ids2 "github.com/isaac-olle/gas/internal/shared/domain/vo/ids"
	securityevents "github.com/isaac-olle/gas/internal/shared/events/security"

	"github.com/isaac-olle/gas/internal/bounded_contexts/backoffice/clients"
	"github.com/isaac-olle/gas/internal/bounded_contexts/backoffice/clients/domain_services"
)

type AddKeyOnClientAssertionKeyCreatedEventHandler struct {
	finder     *domain_services.ClientFinder
	repository clients.Repository
	eventBus   event.Bus
}

func NewAddKeyOnClientAssertionKeyCreatedEventHandler(
	finder *domain_services.ClientFinder,
	repository clients.Repository,
	eventBus event.Bus,
) *AddKeyOnClientAssertionKeyCreatedEventHandler {
	return &AddKeyOnClientAssertionKeyCreatedEventHandler{finder: finder, repository: repository, eventBus: eventBus}
}

func (h *AddKeyOnClientAssertionKeyCreatedEventHandler) Handle(ctx context.Context, evt securityevents.ClientAssertionKeyCreated) error {
	clientId, err := ids2.NewClientIdFromString(evt.ClientId)
	if err != nil {
		return app_errors.BadRequestError(clients.ErrClientNotFound).
			WithDescription("invalid client id in ClientAssertionKeyCreated event").
			WithCause(err)
	}

	client, err := h.finder.Find(ctx, clientId)
	if err != nil {
		return err
	}

	keyId, err := ids2.NewClientAssertionKeyIdFromString(evt.AggregateID())
	if err != nil {
		return app_errors.BadRequestError(clients.ErrClientNotSaved).
			WithDescription("invalid assertion key id in ClientAssertionKeyCreated event").
			WithCause(err)
	}

	client.AddKeyId(keyId)

	if err = h.repository.Save(ctx, client); err != nil {
		return app_errors.InternalError(clients.ErrClientNotSaved).
			WithDescription("failed to save client after adding key id").
			WithCause(err)
	}

	h.eventBus.Publish(ctx, client.PullDomainEvents())
	return nil
}
