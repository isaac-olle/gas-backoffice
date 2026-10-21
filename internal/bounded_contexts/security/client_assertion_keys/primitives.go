package client_assertion_keys

type ClientAssertionKeyPrimitives struct {
	Id             string `json:"id"`
	ClientId       string `json:"client_id"`
	SubscriptionId string `json:"subscription_id"`
	Algorithm      string `json:"algorithm"`
	PublicKey      []byte `json:"public_key"`
}

type OwnershipPrimitives struct {
	ClientId       string `json:"client_id"`
	SubscriptionId string `json:"subscription_id"`
}

func (o Ownership) ToPrimitives() OwnershipPrimitives {
	return OwnershipPrimitives{
		ClientId:       o.ClientId.String(),
		SubscriptionId: o.SubscriptionId.String(),
	}
}

func (c *ClientAssertionKey) ToPrimitives() ClientAssertionKeyPrimitives {
	return ClientAssertionKeyPrimitives{
		Id:             c.id.String(),
		ClientId:       c.ownership.ClientId.String(),
		SubscriptionId: c.ownership.SubscriptionId.String(),
		Algorithm:      c.algorithm.String(),
		PublicKey:      c.key.Bytes(),
	}
}
