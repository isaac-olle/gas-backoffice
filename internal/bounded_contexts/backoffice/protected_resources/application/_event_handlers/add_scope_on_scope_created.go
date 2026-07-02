package _event_handlers

import (
	"context"
	"github.com/isaac-olle/gas/internal/shared/app_errors"
	"github.com/isaac-olle/gas/internal/shared/buses/event"
	ids2 "github.com/isaac-olle/gas/internal/shared/domain/vo/ids"
	backofficeevents "github.com/isaac-olle/gas/internal/shared/events/backoffice"

	"github.com/isaac-olle/gas/internal/bounded_contexts/backoffice/protected_resources"
)

type AddScopeOnScopeCreatedEventHandler struct {
	repository protected_resources.Repository
	eventBus   event.Bus
}

func NewAddScopeOnScopeCreatedEventHandler(repository protected_resources.Repository, eventBus event.Bus) *AddScopeOnScopeCreatedEventHandler {
	return &AddScopeOnScopeCreatedEventHandler{repository: repository, eventBus: eventBus}
}

func (h *AddScopeOnScopeCreatedEventHandler) Handle(ctx context.Context, evt backofficeevents.ScopeCreated) error {
	prId, err := ids2.NewProtectedResourceIdFromString(evt.ProtectedResourceId)
	if err != nil {
		return app_errors.BadRequestError(protected_resources.ErrProtectedResourceNotFound).
			WithDescription("invalid protected resource id in ScopeCreated event").
			WithCause(err)
	}

	pr, err := h.repository.SearchById(ctx, prId)
	if err != nil {
		return app_errors.NotFoundError(protected_resources.ErrProtectedResourceNotFound).
			WithDescription("protected resource not found when handling ScopeCreated event").
			WithCause(err)
	}

	scopeId, err := ids2.NewScopeIdFromString(evt.AggregateID())
	if err != nil {
		return err
	}

	pr.AddScopeId(scopeId)

	if err = h.repository.Save(ctx, pr); err != nil {
		return app_errors.InternalError(protected_resources.ErrProtectedResourceNotSaved).
			WithDescription("failed to save protected resource after adding scope").
			WithCause(err)
	}

	h.eventBus.Publish(ctx, pr.PullDomainEvents())
	return nil
}
