package _event_handlers

import (
	"context"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	backofficeevents "github.com/isaac-olle/gas-backoffice/internal/shared/events/backoffice"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/users"
)

type AddClientOnClientCreatedEventHandler struct {
	repository users.Repository
	eventBus   event.Bus
}

func NewAddClientOnClientCreatedEventHandler(repository users.Repository, eventBus event.Bus) *AddClientOnClientCreatedEventHandler {
	return &AddClientOnClientCreatedEventHandler{repository: repository, eventBus: eventBus}
}

func (h *AddClientOnClientCreatedEventHandler) Handle(ctx context.Context, evt backofficeevents.ClientCreated) error {
	userId, err := ids2.NewUserIdFromString(evt.UserId)
	if err != nil {
		return err
	}

	user, err := h.repository.SearchById(ctx, userId)
	if err != nil {
		return app_errors.NotFoundError(users.ErrUserNotFound).
			WithDescription("user not found when handling ClientCreated event").
			WithCause(err)
	}

	clientId, err := ids2.NewClientIdFromString(evt.AggregateID())
	if err != nil {
		return err
	}

	user.AddClientId(clientId)

	if err = h.repository.Save(ctx, user); err != nil {
		return app_errors.InternalError(users.ErrUserNotSaved).WithCause(err)
	}

	h.eventBus.Publish(ctx, user.PullDomainEvents())
	return nil
}
