package backoffice

import (
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/aggregate"
)

const ScopeDeletedType event.Type = "backoffice.scope.deleted"

type ScopeDeleted struct {
	ProtectedResourceId string
	event.BaseEvent
}

func (e ScopeDeleted) Type() event.Type { return ScopeDeletedType }

func NewScopeDeleted(aggregateId, protectedResourceId string) ScopeDeleted {
	return ScopeDeleted{
		ProtectedResourceId: protectedResourceId,
		BaseEvent: event.NewEvent(
			aggregateId,
			ScopeDeletedType,
			false,
			map[string]any{"protected_resource_id": protectedResourceId},
			aggregate.Scope,
		),
	}
}
