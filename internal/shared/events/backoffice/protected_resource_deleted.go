package backoffice

import (
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/aggregate"
)

const ProtectedResourceDeletedType event.Type = "backoffice.protected_resource.deleted"

type ProtectedResourceDeleted struct {
	UserId string
	event.BaseEvent
}

func (e ProtectedResourceDeleted) Type() event.Type { return ProtectedResourceDeletedType }

func NewProtectedResourceDeleted(aggregateId, userId string) ProtectedResourceDeleted {
	return ProtectedResourceDeleted{
		UserId: userId,
		BaseEvent: event.NewEvent(
			aggregateId,
			ProtectedResourceDeletedType,
			false,
			map[string]any{"user_id": userId},
			aggregate.ProtectedResource,
		),
	}
}
