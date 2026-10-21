package backoffice

import (
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/aggregate"
)

const ClientDeletedType event.Type = "backoffice.client.deleted"

type ClientDeleted struct {
	SubscriptionId string
	event.BaseEvent
}

func (e ClientDeleted) Type() event.Type { return ClientDeletedType }

func NewClientDeleted(aggregateId, subscriptionId string) ClientDeleted {
	return ClientDeleted{
		SubscriptionId: subscriptionId,
		BaseEvent: event.NewEvent(
			aggregateId,
			ClientDeletedType,
			false,
			map[string]any{"subscription_id": subscriptionId},
			aggregate.Client,
		),
	}
}
