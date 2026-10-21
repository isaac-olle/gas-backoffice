package _event_handlers

import (
	"context"
	"errors"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients/domain_services"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	securityevents "github.com/isaac-olle/gas-backoffice/internal/shared/events/security"
	"github.com/isaac-olle/gas-backoffice/internal/shared/persistance"
)

type AddSecretOnClientSecretCreatedEventHandler struct {
	clientFinder *domain_services.ClientFinder
	repository   clients.Repository
}

func NewAddSecretOnClientSecretCreatedEventHandler(
	clientFinder *domain_services.ClientFinder,
	repository clients.Repository,
) *AddSecretOnClientSecretCreatedEventHandler {
	return &AddSecretOnClientSecretCreatedEventHandler{clientFinder: clientFinder, repository: repository}
}

func (h *AddSecretOnClientSecretCreatedEventHandler) Handle(ctx context.Context, evt securityevents.ClientSecretCreated) error {
	clientId, err := ids2.NewClientIdFromString(evt.ClientId)
	if err != nil {
		return app_errors.BadRequestError(clients.ErrClientNotFound).
			WithDescription("invalid client id in ClientSecretCreated event").WithCause(err, false)
	}

	secretId, err := ids2.NewClientSecretIdFromString(evt.AggregateID())
	if err != nil {
		return app_errors.BadRequestError(clients.ErrClientNotFound).
			WithDescription("invalid secret id in ClientSecretCreated event").WithCause(err, false)
	}

	client, err := h.clientFinder.Find(ctx, clientId)
	if err != nil {
		return app_errors.AddTrace(err)
	}

	client.AddSecretId(secretId)

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
