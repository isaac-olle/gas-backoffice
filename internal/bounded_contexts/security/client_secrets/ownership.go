package client_secrets

import (
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

// Ownership is flattened from the owning Client at creation time: ClientId plus
// the Client's own SubscriptionId. It is resolved once (from client_projection)
// and stored locally, so later operations never need to look the Client up again.
type Ownership struct {
	ClientId       ids.ClientId
	SubscriptionId ids.SubscriptionId
}

func NewOwnershipFromVos(clientId ids.ClientId, subscriptionId ids.SubscriptionId) Ownership {
	return Ownership{ClientId: clientId, SubscriptionId: subscriptionId}
}

func NewOwnership(clientId, subscriptionId string) (Ownership, error) {
	cId, err := ids.NewClientIdFromString(clientId)
	if err != nil {
		return Ownership{}, app_errors.BadRequestError(ErrInvalidOwnership).
			WithDescription("invalid client id").WithCause(err, false)
	}
	sId, err := ids.NewSubscriptionIdFromString(subscriptionId)
	if err != nil {
		return Ownership{}, app_errors.BadRequestError(ErrInvalidOwnership).
			WithDescription("invalid subscription id").WithCause(err, false)
	}
	return NewOwnershipFromVos(cId, sId), nil
}
