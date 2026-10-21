package clients

type ClientPrimitives struct {
	Id                string   `json:"id"`
	Name              string   `json:"name"`
	SubscriptionId    string   `json:"subscription_id"`
	Enabled           bool     `json:"enabled"`
	Public            bool     `json:"public"`
	IsSystem          bool     `json:"is_system"`
	AllowedGrantTypes []string `json:"allowed_grant_types"`
	RedirectUris      []string `json:"redirect_uris"`
	AssertionKeyIds   []string `json:"assertion_key_ids"`
	SecretIds         []string `json:"secret_ids"`
	AllowedScopes     []string `json:"allowed_scopes"`
	Version           int      `json:"version"`
}

type OwnershipPrimitives struct {
	SubscriptionId string `json:"subscription_id"`
}

func (o Ownership) ToPrimitives() OwnershipPrimitives {
	return OwnershipPrimitives{
		SubscriptionId: o.SubscriptionId.String(),
	}
}

func (c *Client) ToPrimitives() ClientPrimitives {
	redirectUris := make([]string, 0, c.redirectUris.Size())
	for _, uri := range c.redirectUris.Values() {
		redirectUris = append(redirectUris, uri.String())
	}

	allowedGrantTypes := make([]string, 0, c.allowedGrantTypes.Size())
	for _, g := range c.allowedGrantTypes.Values() {
		allowedGrantTypes = append(allowedGrantTypes, g.String())
	}

	assertionKeyIds := make([]string, 0, c.assertionKeyIds.Size())
	for _, k := range c.assertionKeyIds.Values() {
		assertionKeyIds = append(assertionKeyIds, k.String())
	}

	secretIds := make([]string, 0, c.secretIds.Size())
	for _, s := range c.secretIds.Values() {
		secretIds = append(secretIds, s.String())
	}

	allowedScopes := make([]string, 0, c.allowedScopes.Size())
	for _, sc := range c.allowedScopes.Values() {
		allowedScopes = append(allowedScopes, sc.String())
	}

	return ClientPrimitives{
		Id:                c.id.String(),
		Name:              c.name.String(),
		SubscriptionId:    c.ownership.SubscriptionId.String(),
		Enabled:           c.enabled,
		Public:            c.public,
		IsSystem:          c.isSystem,
		AllowedGrantTypes: allowedGrantTypes,
		RedirectUris:      redirectUris,
		AssertionKeyIds:   assertionKeyIds,
		SecretIds:         secretIds,
		AllowedScopes:     allowedScopes,
		Version:           c.version,
	}
}
