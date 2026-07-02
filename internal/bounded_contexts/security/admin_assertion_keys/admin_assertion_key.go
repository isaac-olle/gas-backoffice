package admin_assertion_keys

import (
	"github.com/isaac-olle/gas/internal/shared/app_errors"
	aggregate2 "github.com/isaac-olle/gas/internal/shared/domain/aggregate"
	vo2 "github.com/isaac-olle/gas/internal/shared/domain/vo"
	ids2 "github.com/isaac-olle/gas/internal/shared/domain/vo/ids"
	securityevents "github.com/isaac-olle/gas/internal/shared/events/security"
)

var adminAssertionKeyAggregateName aggregate2.Name = aggregate2.AdminAssertionKey

type AdminAssertionKey struct {
	*aggregate2.Root
	id        ids2.AdminAssertionKeyId
	ownership Ownership
	algorithm vo2.Algorithm
	key       vo2.PublicKey
}

func NewAdminAssertionKey(id string, adminId ids2.AdminId, algorithm string, publicKey string) (*AdminAssertionKey, error) {
	errArr := app_errors.NewEmptyErrorArray()

	aakId, err := ids2.NewAdminAssertionKeyIdFromString(id)
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

	aak := &AdminAssertionKey{
		Root:      aggregate2.NewRoot(adminAssertionKeyAggregateName),
		id:        aakId,
		ownership: NewOwnershipFromVos(adminId),
		algorithm: alg,
		key:       key,
	}
	aak.Record(securityevents.NewAdminAssertionKeyCreated(aak.id.String(), aak.ownership.AdminId.String()))
	return aak, nil
}

func (aak *AdminAssertionKey) CheckOwnership(o Ownership) bool { return aak.ownership == o }

func FromPrimitives(p AdminAssertionKeyPrimitives) (*AdminAssertionKey, error) {
	errArr := app_errors.NewEmptyErrorArray()

	id, err := ids2.NewAdminAssertionKeyIdFromString(p.Id)
	if err != nil {
		errArr.AddError(err)
	}

	ownership, err := NewOwnership(p.AdminId)
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

	return &AdminAssertionKey{
		Root:      aggregate2.NewRoot(adminAssertionKeyAggregateName),
		id:        id,
		ownership: ownership,
		algorithm: alg,
		key:       key,
	}, nil
}
