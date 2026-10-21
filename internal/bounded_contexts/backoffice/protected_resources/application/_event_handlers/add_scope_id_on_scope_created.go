package _event_handlers

import (
	"context"
	"errors"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/protected_resources"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/protected_resources/domain_services"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	backofficeevents "github.com/isaac-olle/gas-backoffice/internal/shared/events/backoffice"
	"github.com/isaac-olle/gas-backoffice/internal/shared/persistance"
)

type AddScopeIdOnScopeCreatedEventHandler struct {
	prFinder   *domain_services.ProtectedResourceFinder
	repository protected_resources.Repository
}

func NewAddScopeIdOnScopeCreatedEventHandler(
	prFinder *domain_services.ProtectedResourceFinder,
	repository protected_resources.Repository,
) *AddScopeIdOnScopeCreatedEventHandler {
	return &AddScopeIdOnScopeCreatedEventHandler{prFinder: prFinder, repository: repository}
}

func (h *AddScopeIdOnScopeCreatedEventHandler) Handle(ctx context.Context, evt backofficeevents.ScopeCreated) error {
	prId, err := ids2.NewProtectedResourceIdFromString(evt.ProtectedResourceId)
	if err != nil {
		return app_errors.BadRequestError(protected_resources.ErrProtectedResourceNotFound).
			WithDescription("invalid protected resource id in ScopeCreated event").WithCause(err, false)
	}

	scopeId, err := ids2.NewScopeIdFromString(evt.AggregateID())
	if err != nil {
		return app_errors.BadRequestError(protected_resources.ErrProtectedResourceNotFound).
			WithDescription("invalid scope id in ScopeCreated event").WithCause(err, false)
	}

	pr, err := h.prFinder.Find(ctx, prId)
	if err != nil {
		return app_errors.AddTrace(err)
	}

	pr.AddScopeId(scopeId)

	if err = h.repository.Save(ctx, pr); err != nil {
		if errors.Is(err, persistance.ErrConcurrentModification) {
			return app_errors.ConflictError(protected_resources.ErrProtectedResourceConcurrentModification).
				WithDescription("protected resource was modified concurrently; event will be logged as failed").WithCause(err, false)
		}
		return app_errors.InternalError(protected_resources.ErrProtectedResourceNotSaved).
			WithDescription("failed to save protected resource").WithCause(err, false)
	}

	return nil
}
