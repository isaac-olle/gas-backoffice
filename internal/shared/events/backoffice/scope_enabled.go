package backoffice

import (
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/aggregate"
)

const ScopeEnabledType event.Type = "backoffice.scope.enabled"

type ScopeEnabled struct {
	event.BaseEvent
}

func (e ScopeEnabled) Type() event.Type { return ScopeEnabledType }

func NewScopeEnabled(aggregateId string) ScopeEnabled {
	return ScopeEnabled{
		BaseEvent: event.NewEvent(
			aggregateId,
			ScopeEnabledType,
			false,
			map[string]any{},
			aggregate.Scope,
		),
	}
}
