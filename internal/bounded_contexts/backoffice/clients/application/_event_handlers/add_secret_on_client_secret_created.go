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

type AddSecretOnClientSecretCreatedEventHandler struct {
	finder     *domain_services.ClientFinder
	repository clients.Repository
	eventBus   event.Bus
}

func NewAddSecretOnClientSecretCreatedEventHandler(
	finder *domain_services.ClientFinder,
	repository clients.Repository,
	eventBus event.Bus,
) *AddSecretOnClientSecretCreatedEventHandler {
	return &AddSecretOnClientSecretCreatedEventHandler{finder: finder, repository: repository, eventBus: eventBus}
}

func (h *AddSecretOnClientSecretCreatedEventHandler) Handle(ctx context.Context, evt securityevents.ClientSecretCreated) error {
	clientId, err := ids2.NewClientIdFromString(evt.ClientId)
	if err != nil {
		return app_errors.BadRequestError(clients.ErrClientNotFound).
			WithDescription("invalid client id in ClientSecretCreated event").
			WithCause(err)
	}

	client, err := h.finder.Find(ctx, clientId)
	if err != nil {
		return err
	}

	secretId, err := ids2.NewClientSecretIdFromString(evt.AggregateID())
	if err != nil {
		return app_errors.BadRequestError(clients.ErrClientNotSaved).
			WithDescription("invalid secret id in ClientSecretCreated event").
			WithCause(err)
	}

	client.AddSecretId(secretId)

	if err = h.repository.Save(ctx, client); err != nil {
		return app_errors.InternalError(clients.ErrClientNotSaved).
			WithDescription("failed to save client after adding secret id").
			WithCause(err)
	}

	h.eventBus.Publish(ctx, client.PullDomainEvents())
	return nil
}
