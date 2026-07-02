package _event_handlers

import (
	"context"
	"github.com/isaac-olle/gas/internal/shared/app_errors"
	"github.com/isaac-olle/gas/internal/shared/buses/event"
	ids2 "github.com/isaac-olle/gas/internal/shared/domain/vo/ids"
	backofficeevents "github.com/isaac-olle/gas/internal/shared/events/backoffice"

	"github.com/isaac-olle/gas/internal/bounded_contexts/backoffice/tenants"
)

type AddProtectedResourceOnProtectedResourceCreatedEventHandler struct {
	repository tenants.Repository
	eventBus   event.Bus
}

func NewAddProtectedResourceOnProtectedResourceCreatedEventHandler(repository tenants.Repository, eventBus event.Bus) *AddProtectedResourceOnProtectedResourceCreatedEventHandler {
	return &AddProtectedResourceOnProtectedResourceCreatedEventHandler{repository: repository, eventBus: eventBus}
}

func (h *AddProtectedResourceOnProtectedResourceCreatedEventHandler) Handle(ctx context.Context, evt backofficeevents.ProtectedResourceCreated) error {
	tenantId, err := ids2.NewTenantIdFromString(evt.TenantId)
	if err != nil {
		return err
	}

	tenant, err := h.repository.SearchById(ctx, tenantId)
	if err != nil {
		return app_errors.NotFoundError(tenants.ErrTenantNotFound).
			WithDescription("tenant not found when handling ProtectedResourceCreated event").
			WithCause(err)
	}

	prId, err := ids2.NewProtectedResourceIdFromString(evt.AggregateID())
	if err != nil {
		return err
	}

	tenant.AddProtectedResourceId(prId)

	if err = h.repository.Save(ctx, tenant); err != nil {
		return app_errors.InternalError(tenants.ErrTenantNotSaved).WithCause(err)
	}

	h.eventBus.Publish(ctx, tenant.PullDomainEvents())
	return nil
}
