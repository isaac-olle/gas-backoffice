package backoffice

import (
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/aggregate"
)

const ProtectedResourceDeletedType event.Type = "backoffice.protected_resource.deleted"

type ProtectedResourceDeleted struct {
	SubscriptionId string
	event.BaseEvent
}

func (e ProtectedResourceDeleted) Type() event.Type { return ProtectedResourceDeletedType }

func NewProtectedResourceDeleted(aggregateId, subscriptionId string) ProtectedResourceDeleted {
	return ProtectedResourceDeleted{
		SubscriptionId: subscriptionId,
		BaseEvent: event.NewEvent(
			aggregateId,
			ProtectedResourceDeletedType,
			false,
			map[string]any{"subscription_id": subscriptionId},
			aggregate.ProtectedResource,
		),
	}
}
