package _event_handlers

import (
	"context"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	backofficeevents "github.com/isaac-olle/gas-backoffice/internal/shared/events/backoffice"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/users"
)

type RemoveProtectedResourceOnProtectedResourceDeletedEventHandler struct {
	repository users.Repository
	eventBus   event.Bus
}

func NewRemoveProtectedResourceOnProtectedResourceDeletedEventHandler(repository users.Repository, eventBus event.Bus) *RemoveProtectedResourceOnProtectedResourceDeletedEventHandler {
	return &RemoveProtectedResourceOnProtectedResourceDeletedEventHandler{repository: repository, eventBus: eventBus}
}

func (h *RemoveProtectedResourceOnProtectedResourceDeletedEventHandler) Handle(ctx context.Context, evt backofficeevents.ProtectedResourceDeleted) error {
	prId, err := ids2.NewProtectedResourceIdFromString(evt.AggregateID())
	if err != nil {
		return err
	}

	userId, err := ids2.NewUserIdFromString(evt.UserId)
	if err != nil {
		return err
	}

	user, err := h.repository.SearchById(ctx, userId)
	if err != nil {
		return app_errors.NotFoundError(users.ErrUserNotFound).
			WithDescription("user not found when handling ProtectedResourceDeleted event").
			WithCause(err)
	}

	user.RemoveProtectedResourceId(prId)

	if err = h.repository.Save(ctx, user); err != nil {
		return app_errors.InternalError(users.ErrUserNotSaved).WithCause(err)
	}

	h.eventBus.Publish(ctx, user.PullDomainEvents())
	return nil
}
