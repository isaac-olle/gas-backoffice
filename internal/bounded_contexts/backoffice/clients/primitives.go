package clients

type ClientPrimitives struct {
	Id            string
	Name          string
	UserId        string
	AllowedFlows  []string
	AllowedScopes []string
	RedirectUris  []string
	KeyIds        []string
	SecretIds     []string
}

type OwnershipPrimitives struct {
	UserId string
}

func (o Ownership) ToPrimitives() OwnershipPrimitives {
	return OwnershipPrimitives{
		UserId: o.UserId.String(),
	}
}

func (c *Client) ToPrimitives() ClientPrimitives {
	allowedScopes := make([]string, 0, c.allowedScopes.Size())
	for _, id := range c.allowedScopes.Values() {
		allowedScopes = append(allowedScopes, id.String())
	}

	redirectUris := make([]string, 0, c.redirectUris.Size())
	for _, uri := range c.redirectUris.Values() {
		redirectUris = append(redirectUris, uri.String())
	}

	allowedFlows := make([]string, 0, c.allowedFlows.Size())
	for _, f := range c.allowedFlows.Values() {
		allowedFlows = append(allowedFlows, f.String())
	}

	keyIds := make([]string, 0, c.keyIds.Size())
	for _, k := range c.keyIds.Values() {
		keyIds = append(keyIds, k.String())
	}

	secretIds := make([]string, 0, c.secretIds.Size())
	for _, s := range c.secretIds.Values() {
		secretIds = append(secretIds, s.String())
	}

	return ClientPrimitives{
		Id:            c.id.String(),
		Name:          c.name.String(),
		UserId:        c.ownership.UserId.String(),
		AllowedFlows:  allowedFlows,
		AllowedScopes: allowedScopes,
		RedirectUris:  redirectUris,
		KeyIds:        keyIds,
		SecretIds:     secretIds,
	}
}
