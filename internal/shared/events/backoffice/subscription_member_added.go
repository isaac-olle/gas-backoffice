package backoffice

import (
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/aggregate"
)

const SubscriptionMemberAddedType event.Type = "backoffice.subscription.member_added"

type SubscriptionMemberAdded struct {
	UserId           string
	SubscriptionRole string
	event.BaseEvent
}

func (e SubscriptionMemberAdded) Type() event.Type { return SubscriptionMemberAddedType }
func NewSubscriptionMemberAdded(subscriptionId, userId, role string) SubscriptionMemberAdded {
	return SubscriptionMemberAdded{UserId: userId,
		SubscriptionRole: role,
		BaseEvent: event.NewEvent(subscriptionId,
			SubscriptionMemberAddedType,
			false,
			map[string]any{"user_id": userId,
				"subscription_role": role},
			aggregate.Name("backoffice.subscription"))}
}
