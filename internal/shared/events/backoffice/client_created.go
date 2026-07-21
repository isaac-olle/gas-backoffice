package backoffice

import (
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/aggregate"
)

const ClientCreatedType event.Type = "backoffice.client.created"

type ClientCreated struct {
	Name   string
	UserId string
	event.BaseEvent
}

func (e ClientCreated) Type() event.Type { return ClientCreatedType }

func NewClientCreated(aggregateId, name, userId string) ClientCreated {
	return ClientCreated{
		Name:   name,
		UserId: userId,
		BaseEvent: event.NewEvent(
			aggregateId,
			ClientCreatedType,
			false,
			map[string]any{"name": name, "user_id": userId},
			aggregate.Client,
		),
	}
}
