package client_assertion_keys

import (
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	aggregate2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/aggregate"
	vo2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo"
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	securityevents "github.com/isaac-olle/gas-backoffice/internal/shared/events/security"
)

var clientAssertionKeyAggregateName aggregate2.Name = aggregate2.ClientAssertionKey

type ClientAssertionKey struct {
	*aggregate2.Root
	id        ids2.ClientAssertionKeyId
	ownership Ownership
	key       vo2.PublicKey
	algorithm vo2.Algorithm
}

func (c *ClientAssertionKey) EnsureSubscriptionAccess(allowed []ids2.SubscriptionId) error {
	for _, id := range allowed {
		if id == c.ownership.SubscriptionId {
			return nil
		}
	}
	return app_errors.ForbiddenError(ErrInvalidOwnership).
		WithDescription("resource does not belong to an accessible subscription")
}

func NewClientAssertionKeyForSubscription(cakId ids2.ClientAssertionKeyId, clientId ids2.ClientId, subscriptionId ids2.SubscriptionId, algorithmString, publicKey string) (*ClientAssertionKey, error) {
	return newClientAssertionKey(cakId, NewOwnershipFromVos(clientId, subscriptionId), algorithmString, publicKey)
}

func newClientAssertionKey(cakId ids2.ClientAssertionKeyId, ownership Ownership, algorithmString, publicKey string) (*ClientAssertionKey, error) {
	errArr := app_errors.NewEmptyErrorArray()

	key, err := vo2.NewPublicKeyFromHexString(publicKey)
	if err != nil {
		errArr.AddError(err)
	}

	algorithm, err := vo2.NewAlgorithm(algorithmString)
	if err != nil {
		errArr.AddError(err)
	}

	if errArr.HasErrors() {
		return nil, errArr
	}

	cak := &ClientAssertionKey{
		Root:      aggregate2.NewRoot(clientAssertionKeyAggregateName),
		id:        cakId,
		ownership: ownership,
		key:       key,
		algorithm: algorithm,
	}
	cak.Record(securityevents.NewClientAssertionKeyCreated(cak.id.String(), cak.ownership.ClientId.String()))
	return cak, nil
}

func FromPrimitives(p ClientAssertionKeyPrimitives) (*ClientAssertionKey, error) {
	errArr := app_errors.NewEmptyErrorArray()

	id, err := ids2.NewClientAssertionKeyIdFromString(p.Id)
	if err != nil {
		errArr.AddError(err)
	}

	ownership, err := NewOwnership(p.ClientId, p.SubscriptionId)
	if err != nil {
		errArr.AddError(err)
	}

	key, err := vo2.NewPublicKey([]byte(p.PublicKey))
	if err != nil {
		errArr.AddError(err)
	}

	algorithm, err := vo2.NewAlgorithm(p.Algorithm)
	if err != nil {
		errArr.AddError(err)
	}

	if errArr.HasErrors() {
		return nil, errArr
	}

	return &ClientAssertionKey{
		Root:      aggregate2.NewRoot(clientAssertionKeyAggregateName),
		id:        id,
		ownership: ownership,
		key:       key,
		algorithm: algorithm,
	}, nil
}

func (c *ClientAssertionKey) EnsureClientOwnership(clientId ids2.ClientId) error {
	if c.ownership.ClientId != clientId {
		return app_errors.ForbiddenError(ErrInvalidOwnership).
			WithDescription("client assertion key does not belong to the requested client")
	}
	return nil
}
