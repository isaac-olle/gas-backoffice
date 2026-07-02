package key_pairs

type KeyPairPrimitives struct {
	Id                  string
	ProtectedResourceId string
	Algorithm           string
	PrivateKey          []byte
	PublicKey           []byte
}

type OwnershipPrimitives struct {
	ProtectedResourceId string
}

func (o Ownership) ToPrimitives() OwnershipPrimitives {
	return OwnershipPrimitives{ProtectedResourceId: o.ProtectedResourceId.String()}
}

func (kp *KeyPair) ToPrimitives() KeyPairPrimitives {
	return KeyPairPrimitives{
		Id:                  kp.id.String(),
		ProtectedResourceId: kp.ownership.ProtectedResourceId.String(),
		Algorithm:           kp.algorithm.String(),
		PrivateKey:          kp.privateKey.Bytes(),
		PublicKey:           kp.publicKey.Bytes(),
	}
}
