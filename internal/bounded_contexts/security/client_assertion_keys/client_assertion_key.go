package client_assertion_keys

import (
	"github.com/isaac-olle/gas/internal/shared/app_errors"
	aggregate2 "github.com/isaac-olle/gas/internal/shared/domain/aggregate"
	vo2 "github.com/isaac-olle/gas/internal/shared/domain/vo"
	ids2 "github.com/isaac-olle/gas/internal/shared/domain/vo/ids"
	securityevents "github.com/isaac-olle/gas/internal/shared/events/security"
)

var clientAssertionKeyAggregateName aggregate2.Name = aggregate2.ClientAssertionKey

type ClientAssertionKey struct {
	*aggregate2.Root
	id        ids2.ClientAssertionKeyId
	ownership Ownership
	key       vo2.PublicKey
	algorithm vo2.Algorithm
}

func NewClientAssertionKey(cakId ids2.ClientAssertionKeyId, clientId ids2.ClientId, algorithmString string, publicKey string) (*ClientAssertionKey, error) {
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
		ownership: NewOwnershipFromVos(clientId),
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

	ownership, err := NewOwnership(p.ClientId)
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

func (c *ClientAssertionKey) CheckOwnership(o Ownership) bool { return c.ownership == o }
