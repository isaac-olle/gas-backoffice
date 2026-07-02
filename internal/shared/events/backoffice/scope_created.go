package backoffice

import (
	"github.com/isaac-olle/gas/internal/shared/buses/event"
	"github.com/isaac-olle/gas/internal/shared/domain/aggregate"
)

const ScopeCreatedType event.Type = "backoffice.scope.created"

type ScopeCreated struct {
	Name                string
	ProtectedResourceId string
	event.BaseEvent
}

func (e ScopeCreated) Type() event.Type { return ScopeCreatedType }

func NewScopeCreated(aggregateId, name, protectedResourceId string) ScopeCreated {
	return ScopeCreated{
		Name:                name,
		ProtectedResourceId: protectedResourceId,
		BaseEvent: event.NewEvent(
			aggregateId,
			ScopeCreatedType,
			false,
			map[string]any{"name": name, "protected_resource_id": protectedResourceId},
			aggregate.Scope,
		),
	}
}
