package clients

import (
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients/vos"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/aggregate"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"github.com/isaac-olle/gas-backoffice/internal/shared/events/backoffice"
)

var clientAggregateName aggregate.Name = aggregate.Client

type Client struct {
	*aggregate.Root
	id                ids.ClientId
	name              vos.ClientName
	ownership         Ownership
	allowedGrantTypes *vo.SetVo[vos.GrantType]
	redirectUris      *vo.SetVo[vos.RedirectUri]
	assertionKeyIds   *vo.SetVo[ids.ClientAssertionKeyId]
	secretIds         *vo.SetVo[ids.ClientSecretId]
	allowedScopes     *vo.SetVo[ids.ScopeId]
	enabled           bool
	public            bool
	isSystem          bool
	version           int
}

func NewClient(id ids.ClientId, name string, ownership Ownership, redirectUris []string, allowedGrantTypes []string, public bool, isSystem bool) (*Client, error) {
	errArr := app_errors.NewEmptyErrorArray()

	clientName, err := vos.NewClientName(name)
	if err != nil {
		errArr.AddError(err)
	}

	uris := vo.NewEmptySetVo[vos.RedirectUri]()
	for _, raw := range redirectUris {
		uri, uErr := vos.NewRedirectUri(raw)
		if uErr != nil {
			errArr.AddError(uErr)
		} else {
			uris.Add(uri)
		}
	}

	grantTypes := vo.NewEmptySetVo[vos.GrantType]()
	for _, raw := range allowedGrantTypes {
		grantType, gErr := vos.NewGrantType(raw)
		if gErr != nil {
			errArr.AddError(gErr)
		} else {
			grantTypes.Add(grantType)
		}
	}

	if errArr.HasErrors() {
		return nil, errArr
	}

	c := &Client{
		Root:              aggregate.NewRoot(clientAggregateName),
		id:                id,
		name:              clientName,
		ownership:         ownership,
		allowedGrantTypes: grantTypes,
		redirectUris:      uris,
		assertionKeyIds:   vo.NewEmptySetVo[ids.ClientAssertionKeyId](),
		secretIds:         vo.NewEmptySetVo[ids.ClientSecretId](),
		allowedScopes:     vo.NewEmptySetVo[ids.ScopeId](),
		enabled:           true,
		public:            public,
		isSystem:          isSystem,
		version:           1,
	}

	if err = c.checkPublicGrantTypes(); err != nil {
		return nil, err
	}

	c.Record(backoffice.NewClientCreated(c.id.String(), c.name.String(), c.ownership.SubscriptionId.String()))
	return c, nil
}

func FromPrimitives(p ClientPrimitives) (*Client, error) {
	errArr := app_errors.NewEmptyErrorArray()

	id, err := ids.NewClientIdFromString(p.Id)
	if err != nil {
		errArr.AddError(err)
	}

	name, err := vos.NewClientName(p.Name)
	if err != nil {
		errArr.AddError(err)
	}

	var ownership Ownership
	subscriptionId, err := ids.NewSubscriptionIdFromString(p.SubscriptionId)
	if err != nil {
		errArr.AddError(err)
	} else {
		ownership = NewOwnership(subscriptionId)
	}

	allowedGrantTypes := vo.NewEmptySetVo[vos.GrantType]()
	for _, raw := range p.AllowedGrantTypes {
		grantType, gErr := vos.NewGrantType(raw)
		if gErr != nil {
			errArr.AddError(gErr)
		} else {
			allowedGrantTypes.Add(grantType)
		}
	}

	redirectUris := vo.NewEmptySetVo[vos.RedirectUri]()
	for _, raw := range p.RedirectUris {
		uri, uErr := vos.NewRedirectUri(raw)
		if uErr != nil {
			errArr.AddError(uErr)
		} else {
			redirectUris.Add(uri)
		}
	}

	assertionKeyIds := vo.NewEmptySetVo[ids.ClientAssertionKeyId]()
	for _, raw := range p.AssertionKeyIds {
		keyId, kErr := ids.NewClientAssertionKeyIdFromString(raw)
		if kErr != nil {
			errArr.AddError(kErr)
		} else {
			assertionKeyIds.Add(keyId)
		}
	}

	secretIds := vo.NewEmptySetVo[ids.ClientSecretId]()
	for _, raw := range p.SecretIds {
		secretId, sErr := ids.NewClientSecretIdFromString(raw)
		if sErr != nil {
			errArr.AddError(sErr)
		} else {
			secretIds.Add(secretId)
		}
	}

	allowedScopes := vo.NewEmptySetVo[ids.ScopeId]()
	for _, raw := range p.AllowedScopes {
		scopeId, scErr := ids.NewScopeIdFromString(raw)
		if scErr != nil {
			errArr.AddError(scErr)
		} else {
			allowedScopes.Add(scopeId)
		}
	}

	if errArr.HasErrors() {
		return nil, errArr
	}

	// A persisted client always has version >= 1 (the table defaults new rows
	// to 1, and Save always writes a positive version). A zero only ever comes
	// from a hand-built ClientPrimitives{} in a test that didn't set it, so
	// defaulting it to 1 here is safe and never masks a real version.
	version := p.Version
	if version == 0 {
		version = 1
	}

	c := &Client{
		Root:              aggregate.NewRoot(clientAggregateName),
		id:                id,
		name:              name,
		ownership:         ownership,
		allowedGrantTypes: allowedGrantTypes,
		redirectUris:      redirectUris,
		assertionKeyIds:   assertionKeyIds,
		secretIds:         secretIds,
		allowedScopes:     allowedScopes,
		enabled:           p.Enabled,
		public:            p.Public,
		isSystem:          p.IsSystem,
		version:           version,
	}

	return c, nil
}

// checkPublicGrantTypes enforces that a public client (one that can never
// hold a client secret) is never granted client_credentials, since that
// grant type requires authenticating to /token with a secret.
func (c *Client) checkPublicGrantTypes() error {
	if !c.public {
		return nil
	}
	clientCredentials, err := vos.NewGrantType(vos.ClientCredentialsGrantType)
	if err != nil {
		return err
	}
	if c.allowedGrantTypes.Contains(clientCredentials) {
		return app_errors.BadRequestError(ErrPublicClientCannotHaveClientCredentials).
			WithDescription("a public client cannot be granted the client_credentials grant type")
	}
	return nil
}

func (c *Client) IsEnabled() bool        { return c.enabled }
func (c *Client) IsPublic() bool         { return c.public }
func (c *Client) IsSystemResource() bool { return c.isSystem }

func (c *Client) Enable() {
	c.enabled = true
	c.Record(backoffice.NewClientEnabled(c.id.String()))
}

func (c *Client) Disable() {
	c.enabled = false
	c.Record(backoffice.NewClientDisabled(c.id.String()))
}

func (c *Client) Delete() {
	c.Record(backoffice.NewClientDeleted(c.id.String(), c.ownership.SubscriptionId.String()))
}

// EnsureOwnership allows the caller when the client belongs to the given
// subscription. Membership of the caller in that subscription is validated
// earlier, at the subscription path param middleware.
func (c *Client) EnsureOwnership(subscriptionIds []ids.SubscriptionId) error {
	for _, id := range subscriptionIds {
		if id == c.ownership.SubscriptionId {
			return nil
		}
	}
	return app_errors.ForbiddenError(ErrInvalidOwnership).
		WithDescription("client does not belong to the provided subscription")
}

func (c *Client) AddAllowedGrantType(grantType vos.GrantType) {
	c.allowedGrantTypes.Add(grantType)
	c.Record(backoffice.NewClientGrantTypeAdded(c.id.String(), grantType.String()))
}

func (c *Client) RemoveAllowedGrantType(grantType vos.GrantType) {
	c.allowedGrantTypes.Remove(grantType)
	c.Record(backoffice.NewClientGrantTypeRemoved(c.id.String(), grantType.String()))
}

func (c *Client) AddRedirectUri(uri vos.RedirectUri) {
	c.redirectUris.Add(uri)
	c.Record(backoffice.NewClientRedirectUriAdded(c.id.String(), uri.String()))
}

func (c *Client) AddKeyId(keyId ids.ClientAssertionKeyId) {
	c.assertionKeyIds.Add(keyId)
}

func (c *Client) AddSecretId(secretId ids.ClientSecretId) {
	c.secretIds.Add(secretId)
}

func (c *Client) AddAllowedScope(scopeId ids.ScopeId) {
	c.allowedScopes.Add(scopeId)
}

func (c *Client) RemoveAllowedScope(scopeId ids.ScopeId) {
	c.allowedScopes.Remove(scopeId)
}
