package _event_handlers

import (
	"context"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	backofficeevents "github.com/isaac-olle/gas-backoffice/internal/shared/events/backoffice"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/admins"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/admins/domain_services"
)

type RemoveUserOnUserDeletedEventHandler struct {
	repository admins.Repository
	finder     *domain_services.AdminFinder
	eventBus   event.Bus
}

func NewRemoveUserOnUserDeletedEventHandler(repository admins.Repository, finder *domain_services.AdminFinder, eventBus event.Bus) *RemoveUserOnUserDeletedEventHandler {
	return &RemoveUserOnUserDeletedEventHandler{repository: repository, finder: finder, eventBus: eventBus}
}

func (h *RemoveUserOnUserDeletedEventHandler) Handle(ctx context.Context, evt backofficeevents.UserDeleted) error {
	adminId, err := ids2.NewAdminIdFromString(evt.AdminId)
	if err != nil {
		return app_errors.BadRequestError(admins.ErrAdminNotFound).
			WithDescription("invalid admin id in UserDeleted event").
			WithCause(err)
	}

	admin, err := h.finder.Find(ctx, adminId)
	if err != nil {
		return app_errors.AddTrace(err)
	}

	userId, err := ids2.NewUserIdFromString(evt.AggregateID())
	if err != nil {
		return err
	}

	admin.RemoveUserId(userId)

	if err = h.repository.Save(ctx, admin); err != nil {
		return app_errors.InternalError(admins.ErrAdminNotSaved).
			WithDescription("failed to save admin after removing user").
			WithCause(err)
	}

	h.eventBus.Publish(ctx, admin.PullDomainEvents())
	return nil
}
