package client_secrets

import (
	"fmt"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	aggregate2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/aggregate"
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	securityevents "github.com/isaac-olle/gas-backoffice/internal/shared/events/security"
)

var clientSecretAggregateName aggregate2.Name = aggregate2.ClientSecret

type ClientSecret struct {
	*aggregate2.Root
	id         ids2.ClientSecretId
	ownership  Ownership
	secretHash string
}

func NewClientSecretForSubscription(id ids2.ClientSecretId, clientId ids2.ClientId, subscriptionId ids2.SubscriptionId, plainSecret string, hasher SecretHasher) (*ClientSecret, error) {
	return newClientSecret(id, NewOwnershipFromVos(clientId, subscriptionId), plainSecret, hasher)
}

func newClientSecret(id ids2.ClientSecretId, ownership Ownership, plainSecret string, hasher SecretHasher) (*ClientSecret, error) {
	errArr := app_errors.NewEmptyErrorArray()

	if plainSecret == "" {
		errArr.AddError(fmt.Errorf("plain secret cannot be empty"))
	}

	if errArr.HasErrors() {
		return nil, errArr
	}

	secretHash, err := hasher.Hash(plainSecret)
	if err != nil {
		return nil, app_errors.InternalError(ErrClientSecretNotCreated).
			WithDescription("failed to hash client secret").WithCause(
			err, false)
	}

	cs := &ClientSecret{
		Root:       aggregate2.NewRoot(clientSecretAggregateName),
		id:         id,
		ownership:  ownership,
		secretHash: secretHash,
	}
	cs.Record(securityevents.NewClientSecretCreated(cs.id.String(), cs.ownership.ClientId.String()))
	return cs, nil
}

func FromPrimitives(p ClientSecretPrimitives) (*ClientSecret, error) {
	errArr := app_errors.NewEmptyErrorArray()

	id, err := ids2.NewClientSecretIdFromString(p.Id)
	if err != nil {
		errArr.AddError(err)
	}

	ownership, err := NewOwnership(p.ClientId, p.SubscriptionId)
	if err != nil {
		errArr.AddError(err)
	}

	if p.SecretHash == "" {
		errArr.AddError(fmt.Errorf("secret hash cannot be empty"))
	}

	if errArr.HasErrors() {
		return nil, errArr
	}

	return &ClientSecret{
		Root:       aggregate2.NewRoot(clientSecretAggregateName),
		id:         id,
		ownership:  ownership,
		secretHash: p.SecretHash,
	}, nil
}

// EnsureOwnership requires the secret to target clientId and to belong to one
// of subscriptionIds. Membership of the caller in that subscription is
// validated earlier, at the subscription path param middleware.
func (cs *ClientSecret) EnsureOwnership(clientId ids2.ClientId, subscriptionIds []ids2.SubscriptionId) error {
	allowed := false
	for _, id := range subscriptionIds {
		if id == cs.ownership.SubscriptionId {
			allowed = true
			break
		}
	}
	if cs.ownership.ClientId != clientId || !allowed {
		return app_errors.ForbiddenError(ErrInvalidOwnership).
			WithDescription("client secret does not belong to the provided subscription")
	}
	return nil
}
