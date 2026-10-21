package client_secrets

type ClientSecretPrimitives struct {
	Id             string `json:"id"`
	ClientId       string `json:"client_id"`
	SubscriptionId string `json:"subscription_id"`
	SecretHash     string `json:"secret_hash"`
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

func (cs *ClientSecret) ToPrimitives() ClientSecretPrimitives {
	return ClientSecretPrimitives{
		Id:             cs.id.String(),
		ClientId:       cs.ownership.ClientId.String(),
		SubscriptionId: cs.ownership.SubscriptionId.String(),
		SecretHash:     cs.secretHash,
	}
}
