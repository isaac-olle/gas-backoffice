package backoffice

import (
	"github.com/isaac-olle/gas/internal/shared/buses/event"
	"github.com/isaac-olle/gas/internal/shared/domain/aggregate"
)

const ProtectedResourceNameUpdatedType event.Type = "backoffice.protected_resource.name_updated"

type ProtectedResourceNameUpdated struct {
	Name string
	event.BaseEvent
}

func (e ProtectedResourceNameUpdated) Type() event.Type { return ProtectedResourceNameUpdatedType }

func NewProtectedResourceNameUpdated(aggregateId, name string) ProtectedResourceNameUpdated {
	return ProtectedResourceNameUpdated{
		Name: name,
		BaseEvent: event.NewEvent(
			aggregateId,
			ProtectedResourceNameUpdatedType,
			false,
			map[string]any{"name": name},
			aggregate.ProtectedResource,
		),
	}
}
