package backoffice

import (
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/aggregate"
)

const SubscriptionMemberRemovedType event.Type = "backoffice.subscription.member_removed"

type SubscriptionMemberRemoved struct {
	UserId string
	event.BaseEvent
}

func (e SubscriptionMemberRemoved) Type() event.Type { return SubscriptionMemberRemovedType }
func NewSubscriptionMemberRemoved(subscriptionId, userId string) SubscriptionMemberRemoved {
	return SubscriptionMemberRemoved{UserId: userId, BaseEvent: event.NewEvent(subscriptionId, SubscriptionMemberRemovedType, false, map[string]any{"user_id": userId}, aggregate.Name("backoffice.subscription"))}
}
