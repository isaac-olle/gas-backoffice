package _event_handlers

import (
	"context"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/scopes"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/criteria"
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	backofficeevents "github.com/isaac-olle/gas-backoffice/internal/shared/events/backoffice"
)

// DisableScopesOnProtectedResourceDeletedEventHandler disables (does not
// delete) every scope belonging to a protected resource once that protected
// resource is deleted, so they stop being usable without losing their data.
type DisableScopesOnProtectedResourceDeletedEventHandler struct {
	repository scopes.Repository
	eventBus   event.Bus
}

func NewDisableScopesOnProtectedResourceDeletedEventHandler(repository scopes.Repository, eventBus event.Bus) *DisableScopesOnProtectedResourceDeletedEventHandler {
	return &DisableScopesOnProtectedResourceDeletedEventHandler{repository: repository, eventBus: eventBus}
}

func (h *DisableScopesOnProtectedResourceDeletedEventHandler) Handle(ctx context.Context, evt backofficeevents.ProtectedResourceDeleted) error {
	prId, err := ids2.NewProtectedResourceIdFromString(evt.AggregateID())
	if err != nil {
		return app_errors.BadRequestError(scopes.ErrScopeNotFound).
			WithDescription("invalid protected resource id in ProtectedResourceDeleted event").WithCause(err, false)
	}

	affectedScopes, err := h.repository.Matching(ctx, criteria.Criteria{
		Filters: []criteria.Filter{{Field: "protected_resource_id", Operator: criteria.Equals, Value: prId.String()}},
	})
	if err != nil {
		return app_errors.InternalError(scopes.ErrScopeNotSaved).
			WithDescription("failed to search scopes of the deleted protected resource").WithCause(err, false)
	}

	for _, scope := range affectedScopes {
		if !scope.IsEnabled() {
			continue
		}
		scope.Disable()
		if err = h.repository.Save(ctx, scope); err != nil {
			return app_errors.InternalError(scopes.ErrScopeNotSaved).
				WithDescription("failed to disable scope of deleted protected resource").WithCause(err, false)
		}
		h.eventBus.Publish(ctx, scope.PullDomainEvents())
	}

	return nil
}
