package _event_handlers

import (
	"context"
	"github.com/isaac-olle/gas/internal/bounded_contexts/backoffice/admins/domain_services"
	"github.com/isaac-olle/gas/internal/shared/app_errors"
	"github.com/isaac-olle/gas/internal/shared/buses/event"
	ids2 "github.com/isaac-olle/gas/internal/shared/domain/vo/ids"
	backofficeevents "github.com/isaac-olle/gas/internal/shared/events/backoffice"

	"github.com/isaac-olle/gas/internal/bounded_contexts/backoffice/admins"
)

type AddTenantOnTenantCreatedEventHandler struct {
	repository admins.Repository
	finder     *domain_services.AdminFinder
	eventBus   event.Bus
}

func NewAddTenantOnTenantCreatedEventHandler(repository admins.Repository, finder *domain_services.AdminFinder, eventBus event.Bus) *AddTenantOnTenantCreatedEventHandler {
	return &AddTenantOnTenantCreatedEventHandler{repository: repository, finder: finder, eventBus: eventBus}
}

func (h *AddTenantOnTenantCreatedEventHandler) Handle(ctx context.Context, evt backofficeevents.TenantCreated) error {
	adminId, err := ids2.NewAdminIdFromString(evt.AdminId)
	if err != nil {
		return app_errors.BadRequestError(admins.ErrAdminNotFound).
			WithDescription("invalid admin id in TenantCreated event").
			WithCause(err)
	}

	admin, err := h.finder.Find(ctx, adminId)
	if err != nil {
		return app_errors.AddTrace(err)
	}

	tenantId, err := ids2.NewTenantIdFromString(evt.AggregateID())
	if err != nil {
		return app_errors.AddTrace(err)
	}

	admin.AddTenantId(tenantId)

	if err = h.repository.Save(ctx, admin); err != nil {
		return app_errors.InternalError(admins.ErrAdminNotSaved).
			WithDescription("failed to save admin after adding tenant").
			WithCause(err)
	}

	h.eventBus.Publish(ctx, admin.PullDomainEvents())
	return nil
}
