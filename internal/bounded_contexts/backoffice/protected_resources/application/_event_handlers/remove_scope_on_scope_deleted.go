package _event_handlers

import (
	"context"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	backofficeevents "github.com/isaac-olle/gas-backoffice/internal/shared/events/backoffice"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/protected_resources"
)

type RemoveScopeOnScopeDeletedEventHandler struct {
	repository protected_resources.Repository
	eventBus   event.Bus
}

func NewRemoveScopeOnScopeDeletedEventHandler(repository protected_resources.Repository, eventBus event.Bus) *RemoveScopeOnScopeDeletedEventHandler {
	return &RemoveScopeOnScopeDeletedEventHandler{repository: repository, eventBus: eventBus}
}

func (h *RemoveScopeOnScopeDeletedEventHandler) Handle(ctx context.Context, evt backofficeevents.ScopeDeleted) error {
	prId, err := ids2.NewProtectedResourceIdFromString(evt.ProtectedResourceId)
	if err != nil {
		return app_errors.BadRequestError(protected_resources.ErrProtectedResourceNotFound).
			WithDescription("invalid protected resource id in ScopeDeleted event").
			WithCause(err)
	}

	pr, err := h.repository.SearchById(ctx, prId)
	if err != nil {
		return app_errors.NotFoundError(protected_resources.ErrProtectedResourceNotFound).
			WithDescription("protected resource not found when handling ScopeDeleted event").
			WithCause(err)
	}

	scopeId, err := ids2.NewScopeIdFromString(evt.AggregateID())
	if err != nil {
		return err
	}

	pr.RemoveScopeId(scopeId)

	if err = h.repository.Save(ctx, pr); err != nil {
		return app_errors.InternalError(protected_resources.ErrProtectedResourceNotSaved).
			WithDescription("failed to save protected resource after removing scope").
			WithCause(err)
	}

	h.eventBus.Publish(ctx, pr.PullDomainEvents())
	return nil
}
