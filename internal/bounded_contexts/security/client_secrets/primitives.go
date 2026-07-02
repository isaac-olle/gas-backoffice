package client_secrets

type ClientSecretPrimitives struct {
	Id         string
	ClientId   string
	SecretHash string
}

type OwnershipPrimitives struct {
	ClientId string
}

func (o Ownership) ToPrimitives() OwnershipPrimitives {
	return OwnershipPrimitives{ClientId: o.ClientId.String()}
}

func (cs *ClientSecret) ToPrimitives() ClientSecretPrimitives {
	return ClientSecretPrimitives{
		Id:         cs.id.String(),
		ClientId:   cs.ownership.ClientId.String(),
		SecretHash: cs.secretHash,
	}
}
