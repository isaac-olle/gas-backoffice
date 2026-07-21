package user_assertion_keys

import (
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	aggregate2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/aggregate"
	vo2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo"
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	securityevents "github.com/isaac-olle/gas-backoffice/internal/shared/events/security"
)

var userAssertionKeyAggregateName aggregate2.Name = aggregate2.UserAssertionKey

type UserAssertionKey struct {
	*aggregate2.Root
	id        ids2.UserAssertionKeyId
	ownership Ownership
	algorithm vo2.Algorithm
	key       vo2.PublicKey
}

func NewUserAssertionKey(id string, userId ids2.UserId, algorithm string, publicKey string) (*UserAssertionKey, error) {
	errArr := app_errors.NewEmptyErrorArray()

	uakId, err := ids2.NewUserAssertionKeyIdFromString(id)
	if err != nil {
		errArr.AddError(err)
	}

	alg, err := vo2.NewAlgorithm(algorithm)
	if err != nil {
		errArr.AddError(err)
	}

	key, err := vo2.NewPublicKeyFromHexString(publicKey)
	if err != nil {
		errArr.AddError(err)
	}

	if errArr.HasErrors() {
		return nil, errArr
	}

	uak := &UserAssertionKey{
		Root:      aggregate2.NewRoot(userAssertionKeyAggregateName),
		id:        uakId,
		ownership: NewOwnershipFromVos(userId),
		algorithm: alg,
		key:       key,
	}
	uak.Record(securityevents.NewUserAssertionKeyCreated(uak.id.String(), uak.ownership.UserId.String()))
	return uak, nil
}

func (uak *UserAssertionKey) CheckOwnership(o Ownership) bool { return uak.ownership == o }

func FromPrimitives(p UserAssertionKeyPrimitives) (*UserAssertionKey, error) {
	errArr := app_errors.NewEmptyErrorArray()

	id, err := ids2.NewUserAssertionKeyIdFromString(p.Id)
	if err != nil {
		errArr.AddError(err)
	}

	ownership, err := NewOwnership(p.UserId)
	if err != nil {
		errArr.AddError(err)
	}

	alg, err := vo2.NewAlgorithm(p.Algorithm)
	if err != nil {
		errArr.AddError(err)
	}

	key, err := vo2.NewPublicKey(p.PublicKey)
	if err != nil {
		errArr.AddError(err)
	}

	if errArr.HasErrors() {
		return nil, errArr
	}

	return &UserAssertionKey{
		Root:      aggregate2.NewRoot(userAssertionKeyAggregateName),
		id:        id,
		ownership: ownership,
		algorithm: alg,
		key:       key,
	}, nil
}
