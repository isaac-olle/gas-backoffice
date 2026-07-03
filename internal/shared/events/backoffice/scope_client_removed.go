package backoffice

import (
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/aggregate"
)

const ScopeClientRemovedType event.Type = "backoffice.scope.client_removed"

type ScopeClientRemoved struct {
	ClientId string
	event.BaseEvent
}

func (e ScopeClientRemoved) Type() event.Type { return ScopeClientRemovedType }

func NewScopeClientRemoved(aggregateId, clientId string) ScopeClientRemoved {
	return ScopeClientRemoved{
		ClientId: clientId,
		BaseEvent: event.NewEvent(
			aggregateId,
			ScopeClientRemovedType,
			false,
			map[string]any{"client_id": clientId},
			aggregate.Scope,
		),
	}
}
