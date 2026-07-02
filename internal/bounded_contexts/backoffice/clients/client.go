package clients

import (
	"github.com/isaac-olle/gas/internal/bounded_contexts/backoffice/clients/vos"
	"github.com/isaac-olle/gas/internal/shared/app_errors"
	aggregate2 "github.com/isaac-olle/gas/internal/shared/domain/aggregate"
	"github.com/isaac-olle/gas/internal/shared/domain/vo"
	ids2 "github.com/isaac-olle/gas/internal/shared/domain/vo/ids"
	"github.com/isaac-olle/gas/internal/shared/events/backoffice"
)

var clientAggregateName aggregate2.Name = aggregate2.Client

type Client struct {
	*aggregate2.Root
	id            ids2.ClientId
	name          vos.ClientName
	ownership     Ownership
	allowedFlows  *vo.SetVo[vos.Flow]
	allowedScopes *vo.SetVo[ids2.ScopeId]
	redirectUris  *vo.SetVo[vos.RedirectUri]
	keyIds        *vo.SetVo[ids2.ClientAssertionKeyId]
	secretIds     *vo.SetVo[ids2.ClientSecretId]
}

func NewClient(id ids2.ClientId, name string, ownership Ownership, redirectUris []string, allowedFlows []string) (*Client, error) {
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

	flows := vo.NewEmptySetVo[vos.Flow]()
	for _, raw := range allowedFlows {
		flow, fErr := vos.NewFlow(raw)
		if fErr != nil {
			errArr.AddError(fErr)
		} else {
			flows.Add(flow)
		}
	}

	if errArr.HasErrors() {
		return nil, errArr
	}

	c := &Client{
		Root:          aggregate2.NewRoot(clientAggregateName),
		id:            id,
		name:          clientName,
		ownership:     ownership,
		allowedFlows:  flows,
		allowedScopes: vo.NewEmptySetVo[ids2.ScopeId](),
		redirectUris:  uris,
		keyIds:        vo.NewEmptySetVo[ids2.ClientAssertionKeyId](),
		secretIds:     vo.NewEmptySetVo[ids2.ClientSecretId](),
	}
	c.Record(backoffice.NewClientCreated(c.id.String(), c.name.String(), c.ownership.AdminId.String(), c.ownership.TenantId.String()))
	return c, nil
}

func FromPrimitives(p ClientPrimitives) (*Client, error) {
	errArr := app_errors.NewEmptyErrorArray()

	id, err := ids2.NewClientIdFromString(p.Id)
	if err != nil {
		errArr.AddError(err)
	}

	name, err := vos.NewClientName(p.Name)
	if err != nil {
		errArr.AddError(err)
	}

	adminId, err := ids2.NewAdminIdFromString(p.AdminId)
	if err != nil {
		errArr.AddError(err)
	}

	tenantId, err := ids2.NewTenantIdFromString(p.TenantId)
	if err != nil {
		errArr.AddError(err)
	}

	allowedFlows := vo.NewEmptySetVo[vos.Flow]()
	for _, raw := range p.AllowedFlows {
		flow, fErr := vos.NewFlow(raw)
		if fErr != nil {
			errArr.AddError(fErr)
		} else {
			allowedFlows.Add(flow)
		}
	}

	allowedScopes := vo.NewEmptySetVo[ids2.ScopeId]()
	for _, rawId := range p.AllowedScopes {
		scopeId, sErr := ids2.NewScopeIdFromString(rawId)
		if sErr != nil {
			errArr.AddError(sErr)
		} else {
			allowedScopes.Add(scopeId)
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

	keyIds := vo.NewEmptySetVo[ids2.ClientAssertionKeyId]()
	for _, raw := range p.KeyIds {
		keyId, kErr := ids2.NewClientAssertionKeyIdFromString(raw)
		if kErr != nil {
			errArr.AddError(kErr)
		} else {
			keyIds.Add(keyId)
		}
	}

	secretIds := vo.NewEmptySetVo[ids2.ClientSecretId]()
	for _, raw := range p.SecretIds {
		secretId, sErr := ids2.NewClientSecretIdFromString(raw)
		if sErr != nil {
			errArr.AddError(sErr)
		} else {
			secretIds.Add(secretId)
		}
	}

	if errArr.HasErrors() {
		return nil, errArr
	}

	c := &Client{
		Root:          aggregate2.NewRoot(clientAggregateName),
		id:            id,
		name:          name,
		ownership:     NewOwnership(adminId, tenantId),
		allowedFlows:  allowedFlows,
		allowedScopes: allowedScopes,
		redirectUris:  redirectUris,
		keyIds:        keyIds,
		secretIds:     secretIds,
	}
	return c, nil
}

func (c *Client) Delete() {
	c.Record(backoffice.NewClientDeleted(c.id.String(), c.ownership.TenantId.String()))
}

func (c *Client) CheckOwnership(o Ownership) bool {
	return c.ownership == o
}

func (c *Client) AddAllowedFlow(flow vos.Flow) {
	c.allowedFlows.Add(flow)
	c.Record(backoffice.NewClientFlowAdded(c.id.String(), flow.String()))
}

func (c *Client) AddRedirectUri(uri vos.RedirectUri) {
	c.redirectUris.Add(uri)
	c.Record(backoffice.NewClientRedirectUriAdded(c.id.String(), uri.String()))
}

func (c *Client) AddAllowedScope(scopeId ids2.ScopeId) {
	c.allowedScopes.Add(scopeId)
}

func (c *Client) RemoveAllowedScope(scopeId ids2.ScopeId) {
	c.allowedScopes.Remove(scopeId)
}

func (c *Client) AddKeyId(keyId ids2.ClientAssertionKeyId) {
	c.keyIds.Add(keyId)
}

func (c *Client) AddSecretId(secretId ids2.ClientSecretId) {
	c.secretIds.Add(secretId)
}
