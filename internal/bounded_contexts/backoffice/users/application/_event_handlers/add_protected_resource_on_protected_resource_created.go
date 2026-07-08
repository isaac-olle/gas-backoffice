package _event_handlers

import (
	"context"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	backofficeevents "github.com/isaac-olle/gas-backoffice/internal/shared/events/backoffice"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/users"
)

type AddProtectedResourceOnProtectedResourceCreatedEventHandler struct {
	repository users.Repository
	eventBus   event.Bus
}

func NewAddProtectedResourceOnProtectedResourceCreatedEventHandler(repository users.Repository, eventBus event.Bus) *AddProtectedResourceOnProtectedResourceCreatedEventHandler {
	return &AddProtectedResourceOnProtectedResourceCreatedEventHandler{repository: repository, eventBus: eventBus}
}

func (h *AddProtectedResourceOnProtectedResourceCreatedEventHandler) Handle(ctx context.Context, evt backofficeevents.ProtectedResourceCreated) error {
	userId, err := ids2.NewUserIdFromString(evt.UserId)
	if err != nil {
		return err
	}

	user, err := h.repository.SearchById(ctx, userId)
	if err != nil {
		return app_errors.NotFoundError(users.ErrUserNotFound).
			WithDescription("user not found when handling ProtectedResourceCreated event").
			WithCause(err)
	}

	prId, err := ids2.NewProtectedResourceIdFromString(evt.AggregateID())
	if err != nil {
		return err
	}

	user.AddProtectedResourceId(prId)

	if err = h.repository.Save(ctx, user); err != nil {
		return app_errors.InternalError(users.ErrUserNotSaved).WithCause(err)
	}

	h.eventBus.Publish(ctx, user.PullDomainEvents())
	return nil
}
