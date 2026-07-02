package client_assertion_keys

type ClientAssertionKeyPrimitives struct {
	Id        string
	ClientId  string
	Algorithm string
	PublicKey []byte
}

type OwnershipPrimitives struct {
	ClientId string
}

func (o Ownership) ToPrimitives() OwnershipPrimitives {
	return OwnershipPrimitives{ClientId: o.ClientId.String()}
}

func (c *ClientAssertionKey) ToPrimitives() ClientAssertionKeyPrimitives {
	return ClientAssertionKeyPrimitives{
		Id:        c.id.String(),
		ClientId:  c.ownership.ClientId.String(),
		Algorithm: c.algorithm.String(),
		PublicKey: c.key.Bytes(),
	}
}
