package clients

import "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"

type Ownership struct {
	SubscriptionId ids.SubscriptionId
}

func (o Ownership) Equals(other Ownership) bool {
	return o == other
}

func NewOwnership(subscriptionId ids.SubscriptionId) Ownership {
	return Ownership{SubscriptionId: subscriptionId}
}
