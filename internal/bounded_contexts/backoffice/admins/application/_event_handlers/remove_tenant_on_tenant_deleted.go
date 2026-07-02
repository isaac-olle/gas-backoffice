package _event_handlers

import (
	"context"
	"github.com/isaac-olle/gas/internal/shared/app_errors"
	"github.com/isaac-olle/gas/internal/shared/buses/event"
	ids2 "github.com/isaac-olle/gas/internal/shared/domain/vo/ids"
	backofficeevents "github.com/isaac-olle/gas/internal/shared/events/backoffice"

	"github.com/isaac-olle/gas/internal/bounded_contexts/backoffice/admins"
	"github.com/isaac-olle/gas/internal/bounded_contexts/backoffice/admins/domain_services"
)

type RemoveTenantOnTenantDeletedEventHandler struct {
	repository admins.Repository
	finder     *domain_services.AdminFinder
	eventBus   event.Bus
}

func NewRemoveTenantOnTenantDeletedEventHandler(repository admins.Repository, finder *domain_services.AdminFinder, eventBus event.Bus) *RemoveTenantOnTenantDeletedEventHandler {
	return &RemoveTenantOnTenantDeletedEventHandler{repository: repository, finder: finder, eventBus: eventBus}
}

func (h *RemoveTenantOnTenantDeletedEventHandler) Handle(ctx context.Context, evt backofficeevents.TenantDeleted) error {
	adminId, err := ids2.NewAdminIdFromString(evt.AdminId)
	if err != nil {
		return app_errors.BadRequestError(admins.ErrAdminNotFound).
			WithDescription("invalid admin id in TenantDeleted event").
			WithCause(err)
	}

	admin, err := h.finder.Find(ctx, adminId)
	if err != nil {
		return app_errors.AddTrace(err)
	}

	tenantId, err := ids2.NewTenantIdFromString(evt.AggregateID())
	if err != nil {
		return err
	}

	admin.RemoveTenantId(tenantId)

	if err = h.repository.Save(ctx, admin); err != nil {
		return app_errors.InternalError(admins.ErrAdminNotSaved).
			WithDescription("failed to save admin after removing tenant").
			WithCause(err)
	}

	h.eventBus.Publish(ctx, admin.PullDomainEvents())
	return nil
}
