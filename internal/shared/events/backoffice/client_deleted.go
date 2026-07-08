package backoffice

import (
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/aggregate"
)

const ClientDeletedType event.Type = "backoffice.client.deleted"

type ClientDeleted struct {
	UserId string
	event.BaseEvent
}

func (e ClientDeleted) Type() event.Type { return ClientDeletedType }

func NewClientDeleted(aggregateId, userId string) ClientDeleted {
	return ClientDeleted{
		UserId: userId,
		BaseEvent: event.NewEvent(
			aggregateId,
			ClientDeletedType,
			false,
			map[string]any{"user_id": userId},
			aggregate.Client,
		),
	}
}
