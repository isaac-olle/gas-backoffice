package scopes

import (
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

type Ownership struct {
	ProtectedResourceId ids2.ProtectedResourceId
	SubscriptionId      ids2.SubscriptionId
}

func (o Ownership) Equals(other Ownership) bool {
	return o == other
}

func NewOwnershipFromVos(protectedResourceId ids2.ProtectedResourceId, subscriptionId ids2.SubscriptionId) Ownership {
	return Ownership{ProtectedResourceId: protectedResourceId, SubscriptionId: subscriptionId}
}

func NewOwnership(protectedResourceId, subscriptionId string) (Ownership, error) {
	prId, err := ids2.NewProtectedResourceIdFromString(protectedResourceId)
	if err != nil {
		return Ownership{}, app_errors.BadRequestError(ErrInvalidOwnership).
			WithDescription("invalid protected resource id").WithCause(err, false)
	}
	sId, err := ids2.NewSubscriptionIdFromString(subscriptionId)
	if err != nil {
		return Ownership{}, app_errors.BadRequestError(ErrInvalidOwnership).
			WithDescription("invalid subscription id").WithCause(err, false)
	}
	return NewOwnershipFromVos(prId, sId), nil
}
