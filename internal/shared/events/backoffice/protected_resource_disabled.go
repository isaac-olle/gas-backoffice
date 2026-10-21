package backoffice

import (
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/aggregate"
)

const ProtectedResourceDisabledType event.Type = "backoffice.protected_resource.disabled"

type ProtectedResourceDisabled struct {
	event.BaseEvent
}

func (e ProtectedResourceDisabled) Type() event.Type { return ProtectedResourceDisabledType }

func NewProtectedResourceDisabled(aggregateId string) ProtectedResourceDisabled {
	return ProtectedResourceDisabled{
		BaseEvent: event.NewEvent(
			aggregateId,
			ProtectedResourceDisabledType,
			false,
			map[string]any{},
			aggregate.ProtectedResource,
		),
	}
}
