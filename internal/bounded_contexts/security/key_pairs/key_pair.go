package key_pairs

import (
	"github.com/isaac-olle/gas/internal/shared/app_errors"
	aggregate2 "github.com/isaac-olle/gas/internal/shared/domain/aggregate"
	vo2 "github.com/isaac-olle/gas/internal/shared/domain/vo"
	"github.com/isaac-olle/gas/internal/shared/domain/vo/ids"
	securityevents "github.com/isaac-olle/gas/internal/shared/events/security"
)

var keyPairAggregateName aggregate2.Name = aggregate2.KeyPair

type KeyPair struct {
	*aggregate2.Root
	id         ids.KeyPairId
	ownership  Ownership
	algorithm  vo2.Algorithm
	privateKey vo2.PrivateKey
	publicKey  vo2.PublicKey
}

func NewKeyPair(id string, protectedResourceId string, algorithm string, privateKeyHex string, publicKeyHex string) (*KeyPair, error) {
	errArr := app_errors.NewEmptyErrorArray()

	kpId, err := ids.NewKeyPairIdFromString(id)
	if err != nil {
		errArr.AddError(err)
	}

	ownership, err := NewOwnership(protectedResourceId)
	if err != nil {
		errArr.AddError(err)
	}

	alg, err := vo2.NewAlgorithm(algorithm)
	if err != nil {
		errArr.AddError(err)
	}

	privKey, err := vo2.NewPrivateKeyFromHexString(privateKeyHex)
	if err != nil {
		errArr.AddError(err)
	}

	pubKey, err := vo2.NewPublicKeyFromHexString(publicKeyHex)
	if err != nil {
		errArr.AddError(err)
	}

	if errArr.HasErrors() {
		return nil, errArr
	}

	kp := &KeyPair{
		Root:       aggregate2.NewRoot(keyPairAggregateName),
		id:         kpId,
		ownership:  ownership,
		algorithm:  alg,
		privateKey: privKey,
		publicKey:  pubKey,
	}
	kp.Record(securityevents.NewKeyPairCreated(kp.id.String(), algorithm))
	return kp, nil
}

func FromPrimitives(p KeyPairPrimitives) (*KeyPair, error) {
	errArr := app_errors.NewEmptyErrorArray()

	id, err := ids.NewKeyPairIdFromString(p.Id)
	if err != nil {
		errArr.AddError(err)
	}

	ownership, err := NewOwnership(p.ProtectedResourceId)
	if err != nil {
		errArr.AddError(err)
	}

	alg, err := vo2.NewAlgorithm(p.Algorithm)
	if err != nil {
		errArr.AddError(err)
	}

	privKey, err := vo2.NewPrivateKey([]byte(p.PrivateKey))
	if err != nil {
		errArr.AddError(err)
	}

	pubKey, err := vo2.NewPublicKey([]byte(p.PublicKey))
	if err != nil {
		errArr.AddError(err)
	}

	if errArr.HasErrors() {
		return nil, errArr
	}

	return &KeyPair{
		Root:       aggregate2.NewRoot(keyPairAggregateName),
		id:         id,
		ownership:  ownership,
		algorithm:  alg,
		privateKey: privKey,
		publicKey:  pubKey,
	}, nil
}

func (kp *KeyPair) Algorithm() vo2.Algorithm        { return kp.algorithm }
func (kp *KeyPair) PrivateKey() vo2.PrivateKey      { return kp.privateKey }
func (kp *KeyPair) PublicKey() vo2.PublicKey        { return kp.publicKey }
func (kp *KeyPair) CheckOwnership(o Ownership) bool { return kp.ownership == o }
