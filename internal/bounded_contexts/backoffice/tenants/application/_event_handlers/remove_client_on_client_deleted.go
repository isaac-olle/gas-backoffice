package _event_handlers

import (
	"context"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	backofficeevents "github.com/isaac-olle/gas-backoffice/internal/shared/events/backoffice"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/tenants"
)

type RemoveClientOnClientDeletedEventHandler struct {
	repository tenants.Repository
	eventBus   event.Bus
}

func NewRemoveClientOnClientDeletedEventHandler(repository tenants.Repository, eventBus event.Bus) *RemoveClientOnClientDeletedEventHandler {
	return &RemoveClientOnClientDeletedEventHandler{repository: repository, eventBus: eventBus}
}

func (h *RemoveClientOnClientDeletedEventHandler) Handle(ctx context.Context, evt backofficeevents.ClientDeleted) error {
	clientId, err := ids2.NewClientIdFromString(evt.AggregateID())
	if err != nil {
		return err
	}

	tenantId, err := ids2.NewTenantIdFromString(evt.TenantId)
	if err != nil {
		return err
	}

	tenant, err := h.repository.SearchById(ctx, tenantId)
	if err != nil {
		return app_errors.NotFoundError(tenants.ErrTenantNotFound).
			WithDescription("tenant not found when handling ClientDeleted event").
			WithCause(err)
	}

	tenant.RemoveClientId(clientId)

	if err = h.repository.Save(ctx, tenant); err != nil {
		return app_errors.InternalError(tenants.ErrTenantNotSaved).WithCause(err)
	}

	h.eventBus.Publish(ctx, tenant.PullDomainEvents())
	return nil
}
