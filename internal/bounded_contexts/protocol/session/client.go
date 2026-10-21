package session

import (
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"strings"

	sessionvos "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session/vos"
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
	ClientID          string   `bson:"client_id"`
	Name              string   `bson:"name"`
	UserId            string   `bson:"user_id"`
	AllowedGrantTypes []string `bson:"allowed_grant_types"`
	RedirectURIs      []string `bson:"redirect_uris"`
	Scopes            []*Scope `bson:"scopes"`
	Public            bool     `bson:"public"`
}

func NewClient(clientID, name, userId string, allowedGrantTypes, redirectURIs []string, scopes []*Scope, public bool) *Client {
	return &Client{
		ClientID:          clientID,
		Name:              name,
		UserId:            userId,
		AllowedGrantTypes: allowedGrantTypes,
		RedirectURIs:      redirectURIs,
		Scopes:            scopes,
		Public:            public,
	}
}

func (c *Client) AssertPublic(publicOnly bool) error {
	if publicOnly && !c.Public {
		return app_errors.ForbiddenError(ErrClientNotPublic).
			WithDescription("this operation only accepts public clients")
	}
	return nil
}

func (c *Client) HasGrantType(grantType string) bool {
	for _, g := range c.AllowedGrantTypes {
		if g == grantType {
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

func (c *Client) CheckClientCredentials(scope string) (string, error) {
	if !c.HasGrantType(sessionvos.ClientCredentialsGrantType.String()) {
		return "", app_errors.BadRequestError(ErrGrantTypeNotAllowed).
			WithDescription("client_credentials grant type is not allowed for this client")
	}
	requestedScopes := strings.Fields(scope)
	for _, scopeName := range requestedScopes {
		if !c.HasScope(scopeName) {
			return "", app_errors.BadRequestError(ErrInvalidScope).
				WithDescription("scope not allowed for this client: " + scopeName)
		}
	}
	return c.ResolveProtectedResourceId(requestedScopes)
}

// ResolveProtectedResourceId validates that all requested scopes are unambiguous and belong
// to the same ProtectedResource, then returns that ProtectedResource ID.
// A scope is ambiguous if the same name exists under multiple ProtectedResources — in that
// case the client must specify the `resource` parameter (RFC 8707).
func (c *Client) ResolveProtectedResourceId(requestedScopes []string) (string, error) {
	var resolvedPrId string

	for _, scopeName := range requestedScopes {
		seen := map[string]string{} // protectedResourceId -> protectedResourceName
		for _, s := range c.Scopes {
			if s.ScopeName == scopeName {
				seen[s.ProtectedResourceID] = s.ProtectedResourceName
			}
		}
		if len(seen) > 1 {
			resources := make([]string, 0, len(seen))
			for _, name := range seen {
				resources = append(resources, name)
			}
			desc := "scope \"" + scopeName + "\" exists in multiple protected resources: " +
				strings.Join(resources, ", ") + ". Use the `resource` parameter (RFC 8707) to specify which one."
			return "", app_errors.BadRequestError(ErrAmbiguousScope).WithDescription(desc)
		}
		for prId := range seen {
			if resolvedPrId == "" {
				resolvedPrId = prId
			} else if resolvedPrId != prId {
				return "", app_errors.BadRequestError(ErrScopesFromMultipleResources).
					WithDescription("requested scopes belong to different protected resources. Request scopes from one resource at a time.")
			}
		}
	}

	return resolvedPrId, nil
}
