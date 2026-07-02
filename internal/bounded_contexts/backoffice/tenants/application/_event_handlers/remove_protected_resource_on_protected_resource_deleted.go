package _event_handlers

import (
	"context"
	"github.com/isaac-olle/gas/internal/shared/app_errors"
	"github.com/isaac-olle/gas/internal/shared/buses/event"
	ids2 "github.com/isaac-olle/gas/internal/shared/domain/vo/ids"
	backofficeevents "github.com/isaac-olle/gas/internal/shared/events/backoffice"

	"github.com/isaac-olle/gas/internal/bounded_contexts/backoffice/tenants"
)

type RemoveProtectedResourceOnProtectedResourceDeletedEventHandler struct {
	repository tenants.Repository
	eventBus   event.Bus
}

func NewRemoveProtectedResourceOnProtectedResourceDeletedEventHandler(repository tenants.Repository, eventBus event.Bus) *RemoveProtectedResourceOnProtectedResourceDeletedEventHandler {
	return &RemoveProtectedResourceOnProtectedResourceDeletedEventHandler{repository: repository, eventBus: eventBus}
}

func (h *RemoveProtectedResourceOnProtectedResourceDeletedEventHandler) Handle(ctx context.Context, evt backofficeevents.ProtectedResourceDeleted) error {
	prId, err := ids2.NewProtectedResourceIdFromString(evt.AggregateID())
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
			WithDescription("tenant not found when handling ProtectedResourceDeleted event").
			WithCause(err)
	}

	tenant.RemoveProtectedResourceId(prId)

	if err = h.repository.Save(ctx, tenant); err != nil {
		return app_errors.InternalError(tenants.ErrTenantNotSaved).WithCause(err)
	}

	h.eventBus.Publish(ctx, tenant.PullDomainEvents())
	return nil
}
