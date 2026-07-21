package user_assertion_keys

type UserAssertionKeyPrimitives struct {
	Id        string
	UserId    string
	Algorithm string
	PublicKey []byte
}

type OwnershipPrimitives struct {
	UserId string
}

func (o Ownership) ToPrimitives() OwnershipPrimitives {
	return OwnershipPrimitives{UserId: o.UserId.String()}
}

func (uak *UserAssertionKey) ToPrimitives() UserAssertionKeyPrimitives {
	return UserAssertionKeyPrimitives{
		Id:        uak.id.String(),
		UserId:    uak.ownership.UserId.String(),
		Algorithm: uak.algorithm.String(),
		PublicKey: uak.key.Bytes(),
	}
}
