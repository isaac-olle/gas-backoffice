package session

import "time"

type BaseSessionPrimitives struct {
	ClientID            string
	Scope               string
	GrantType           string
	ProtectedResourceId string
}

type AuthorizationCodeSessionPrimitives struct {
	*BaseSessionPrimitives
	SessionID         string
	ClientOwnerUserId string
	RedirectURI       string
	State             string
	ResponseTypes     []string
	Code              string
	Pkce              *PkcePrimitives
	SessionState      string
	Consumed          bool
	Identity          UserIdentityPrimitives
}

// UserIdentityPrimitives is the flattened persistence shape of
// sessionvos.ExternalOidcUserIdentity, stored on the same row as the session.
type UserIdentityPrimitives struct {
	IdpId        string
	IDToken      string
	AccessToken  string
	RefreshToken *string
	Claims       IdpClaimsPrimitives
	ExpiresAt    time.Time
}

type IdpClaimsPrimitives struct {
	Sub     string `json:"sub"`
	Email   string `json:"email"`
	Picture string `json:"picture,omitempty"`
}

type PkcePrimitives struct {
	CodeChallenge string
	Method        string
}

type ClientCredentialsSessionPrimitives struct {
	*BaseSessionPrimitives
	ID       string
	UserId   string
	IssuedAt string
}
