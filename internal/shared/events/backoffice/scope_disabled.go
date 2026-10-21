package backoffice

import (
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/aggregate"
)

const ScopeDisabledType event.Type = "backoffice.scope.disabled"

type ScopeDisabled struct {
	event.BaseEvent
}

func (e ScopeDisabled) Type() event.Type { return ScopeDisabledType }

func NewScopeDisabled(aggregateId string) ScopeDisabled {
	return ScopeDisabled{
		BaseEvent: event.NewEvent(
			aggregateId,
			ScopeDisabledType,
			false,
			map[string]any{},
			aggregate.Scope,
		),
	}
}
