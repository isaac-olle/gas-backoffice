package session

import (
	"slices"
	"strings"
)

type TokenClaims struct {
	// Subject is the external OIDC subject scoped to IdpID for authorization
	// code, or the client ID for client_credentials.
	Subject             string
	GrantType           string
	IdpID               string
	Email               string
	Picture             string
	Scope               string
	ClientID            string
	ProtectedResourceId string
}

func (c TokenClaims) HasScope(scope string) bool {
	return slices.Contains(strings.Fields(c.Scope), scope)
}

type Session interface {
	TokenClaims() TokenClaims
}

type AccessTokenClaims struct {
	Issuer string
	Scope  string
}

type IDTokenClaims struct {
	Issuer   string
	Subject  string
	Audience string
}
