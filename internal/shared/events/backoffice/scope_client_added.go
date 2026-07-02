package backoffice

import (
	"github.com/isaac-olle/gas/internal/shared/buses/event"
	"github.com/isaac-olle/gas/internal/shared/domain/aggregate"
)

const ScopeClientAddedType event.Type = "backoffice.scope.client_added"

type ScopeClientAdded struct {
	ClientId string
	event.BaseEvent
}

func (e ScopeClientAdded) Type() event.Type { return ScopeClientAddedType }

func NewScopeClientAdded(aggregateId, clientId string) ScopeClientAdded {
	return ScopeClientAdded{
		ClientId: clientId,
		BaseEvent: event.NewEvent(
			aggregateId,
			ScopeClientAddedType,
			false,
			map[string]any{"client_id": clientId},
			aggregate.Scope,
		),
	}
}
