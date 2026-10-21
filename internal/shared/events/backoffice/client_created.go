package backoffice

import (
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/aggregate"
)

const ClientCreatedType event.Type = "backoffice.client.created"

type ClientCreated struct {
	Name           string
	SubscriptionId string
	event.BaseEvent
}

func (e ClientCreated) Type() event.Type { return ClientCreatedType }

func NewClientCreated(aggregateId, name, subscriptionId string) ClientCreated {
	return ClientCreated{
		Name:           name,
		SubscriptionId: subscriptionId,
		BaseEvent: event.NewEvent(
			aggregateId,
			ClientCreatedType,
			false,
			map[string]any{"name": name, "subscription_id": subscriptionId},
			aggregate.Client,
		),
	}
}
