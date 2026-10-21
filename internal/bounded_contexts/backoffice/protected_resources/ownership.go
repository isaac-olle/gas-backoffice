package protected_resources

import (
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

type Ownership struct {
	SubscriptionId ids2.SubscriptionId
}

func NewOwnership(subscriptionId ids2.SubscriptionId) Ownership {
	return Ownership{SubscriptionId: subscriptionId}
}

func NewOwnershipFromStrings(subscriptionId string) (Ownership, error) {
	sId, err := ids2.NewSubscriptionIdFromString(subscriptionId)
	if err != nil {
		return Ownership{}, app_errors.BadRequestError(ErrInvalidOwnership).
			WithDescription("invalid subscription id").WithCause(err, false)
	}
	return NewOwnership(sId), nil
}

func (o Ownership) Equals(other Ownership) bool {
	return o == other
}
