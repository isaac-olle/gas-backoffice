package _event_handlers

import (
	"context"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/admins/domain_services"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	backofficeevents "github.com/isaac-olle/gas-backoffice/internal/shared/events/backoffice"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/admins"
)

type AddUserOnUserCreatedEventHandler struct {
	repository admins.Repository
	finder     *domain_services.AdminFinder
	eventBus   event.Bus
}

func NewAddUserOnUserCreatedEventHandler(repository admins.Repository, finder *domain_services.AdminFinder, eventBus event.Bus) *AddUserOnUserCreatedEventHandler {
	return &AddUserOnUserCreatedEventHandler{repository: repository, finder: finder, eventBus: eventBus}
}

func (h *AddUserOnUserCreatedEventHandler) Handle(ctx context.Context, evt backofficeevents.UserCreated) error {
	adminId, err := ids2.NewAdminIdFromString(evt.AdminId)
	if err != nil {
		return app_errors.BadRequestError(admins.ErrAdminNotFound).
			WithDescription("invalid admin id in UserCreated event").
			WithCause(err)
	}

	admin, err := h.finder.Find(ctx, adminId)
	if err != nil {
		return app_errors.AddTrace(err)
	}

	userId, err := ids2.NewUserIdFromString(evt.AggregateID())
	if err != nil {
		return app_errors.AddTrace(err)
	}

	admin.AddUserId(userId)

	if err = h.repository.Save(ctx, admin); err != nil {
		return app_errors.InternalError(admins.ErrAdminNotSaved).
			WithDescription("failed to save admin after adding user").
			WithCause(err)
	}

	h.eventBus.Publish(ctx, admin.PullDomainEvents())
	return nil
}
