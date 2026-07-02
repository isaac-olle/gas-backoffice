package session

import (
	"github.com/isaac-olle/gas/internal/shared/app_errors"
	"strings"

	sessionvos "github.com/isaac-olle/gas/internal/bounded_contexts/protocol/session/vos"
)

type Scope struct {
	ScopeID               string `bson:"scope_id"`
	ScopeName             string `bson:"scope_name"`
	ProtectedResourceID   string `bson:"protected_resource_id"`
	ProtectedResourceName string `bson:"protected_resource_name"`
}

func NewScope(scopeID, scopeName, protectedResourceID, protectedResourceName string) *Scope {
	return &Scope{
		ScopeID:               scopeID,
		ScopeName:             scopeName,
		ProtectedResourceID:   protectedResourceID,
		ProtectedResourceName: protectedResourceName,
	}
}

type Client struct {
	ClientID     string   `bson:"client_id"`
	Name         string   `bson:"name"`
	TenantID     string   `bson:"tenant_id"`
	AdminID      string   `bson:"admin_id"`
	AllowedFlows []string `bson:"allowed_flows"`
	RedirectURIs []string `bson:"redirect_uris"`
	Scopes       []*Scope `bson:"scopes"`
}

func NewClient(clientID, name, tenantID, adminID string, allowedFlows, redirectURIs []string, scopes []*Scope) *Client {
	return &Client{
		ClientID:     clientID,
		Name:         name,
		TenantID:     tenantID,
		AdminID:      adminID,
		AllowedFlows: allowedFlows,
		RedirectURIs: redirectURIs,
		Scopes:       scopes,
	}
}

func (c *Client) HasFlow(flow string) bool {
	for _, f := range c.AllowedFlows {
		if f == flow {
			return true
		}
	}
	return false
}

func (c *Client) HasRedirectURI(uri string) bool {
	for _, redirectURI := range c.RedirectURIs {
		if redirectURI == uri {
			return true
		}
	}
	return false
}

func (c *Client) HasScope(scopeName string) bool {
	for _, scope := range c.Scopes {
		if scope.ScopeName == scopeName {
			return true
		}
	}
	return false
}

func (c *Client) CheckAuthorizationRequest(responseType, redirectURI, scope string) error {
	flow, err := sessionvos.DetermineFlow(responseType)
	if err != nil {
		return app_errors.BadRequestError(ErrFlowNotAllowed).
			WithDescription("unsupported response_type: " + responseType)
	}
	if !c.HasFlow(flow.String()) {
		return app_errors.BadRequestError(ErrFlowNotAllowed).
			WithDescription("flow not allowed for this client: " + flow.String())
	}
	if !c.HasRedirectURI(redirectURI) {
		return app_errors.BadRequestError(ErrInvalidRedirectURI).
			WithDescription("redirect_uri not registered for this client")
	}
	for _, scopeName := range strings.Fields(scope) {
		if !c.HasScope(scopeName) {
			return app_errors.BadRequestError(ErrInvalidScope).
				WithDescription("scope not allowed for this client: " + scopeName)
		}
	}
	return nil
}

func (c *Client) CheckClientCredentials(scope string) error {
	if !c.HasFlow(sessionvos.ClientCredentialsFlow) {
		return app_errors.BadRequestError(ErrFlowNotAllowed).
			WithDescription("client_credentials flow is not allowed for this client")
	}
	for _, scopeName := range strings.Fields(scope) {
		if !c.HasScope(scopeName) {
			return app_errors.BadRequestError(ErrInvalidScope).
				WithDescription("scope not allowed for this client: " + scopeName)
		}
	}
	return nil
}
