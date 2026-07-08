package backoffice

import (
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/aggregate"
)

const ProtectedResourceCreatedType event.Type = "backoffice.protected_resource.created"

type ProtectedResourceCreated struct {
	Name   string
	UserId string
	event.BaseEvent
}

func (e ProtectedResourceCreated) Type() event.Type { return ProtectedResourceCreatedType }

func NewProtectedResourceCreated(aggregateId, name, userId string) ProtectedResourceCreated {
	return ProtectedResourceCreated{
		Name:   name,
		UserId: userId,
		BaseEvent: event.NewEvent(
			aggregateId,
			ProtectedResourceCreatedType,
			false,
			map[string]any{"name": name, "user_id": userId},
			aggregate.ProtectedResource,
		),
	}
}
