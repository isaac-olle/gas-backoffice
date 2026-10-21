package backoffice

import (
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/aggregate"
)

const ProtectedResourceEnabledType event.Type = "backoffice.protected_resource.enabled"

type ProtectedResourceEnabled struct {
	event.BaseEvent
}

func (e ProtectedResourceEnabled) Type() event.Type { return ProtectedResourceEnabledType }

func NewProtectedResourceEnabled(aggregateId string) ProtectedResourceEnabled {
	return ProtectedResourceEnabled{
		BaseEvent: event.NewEvent(
			aggregateId,
			ProtectedResourceEnabledType,
			false,
			map[string]any{},
			aggregate.ProtectedResource,
		),
	}
}
