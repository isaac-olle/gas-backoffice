package session

import (
	"context"
	"time"

	protocolShared "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/_shared"
	sessionvos "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session/vos"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

type Repository interface {
	Save(ctx context.Context, sess *AuthorizationCodeSession) error
	FindById(ctx context.Context, id ids.SessionId) (*AuthorizationCodeSession, error)
	FindByCode(ctx context.Context, code string) (*AuthorizationCodeSession, error)
	SoftDelete(ctx context.Context, id ids.SessionId) error
	HardDelete(ctx context.Context, id ids.SessionId) error
}

type ClientFinderRepository interface {
	FindById(ctx context.Context, clientID string) (*Client, error)
}

type AuthorizationURLBuilder interface {
	BuildAuthorizationURL(idp *protocolShared.IdentityProviderConfig, state string) (string, error)
}

// IdpStatePayload is the self-contained, encrypted "state" sent to the
// identity provider at /oidc/authorize and decoded back at the callback by
// create_session. It carries everything /oidc/authorize validated so nothing
// needs to be persisted until the user actually completes the IDP login.
type IdpStatePayload struct {
	ClientID            string `json:"cid"`
	RedirectURI         string `json:"ru"`
	Scope               string `json:"sc"`
	ResponseType        string `json:"rt"`
	CodeChallenge       string `json:"cc"`
	CodeChallengeMethod string `json:"ccm"`
	ClientState         string `json:"cs"`
	IdpReferenceName    string `json:"idp"`
	Exp                 int64  `json:"exp"`
}

// IdpStateCodec encodes/decodes the opaque "state" value sent to and received
// back from the identity provider.
type IdpStateCodec interface {
	Encode(payload IdpStatePayload) (string, error)
	Decode(encoded string) (IdpStatePayload, error)
}

// IdpTokens holds the tokens obtained from exchanging an authorization code
// with an identity provider.
type IdpTokens struct {
	IDToken      string
	AccessToken  string
	RefreshToken string
	ExpiresAt    time.Time
}

// IdpAuthentication holds the result of authenticating a user through an
// identity provider: the resolved identity and the tokens obtained from it.
type IdpAuthentication struct {
	Tokens IdpTokens
	Claims sessionvos.IdpClaims
}

// IdentityProviderClient is the infrastructure adapter that exchanges an
// authorization code for tokens with an external identity provider.
type IdentityProviderClient interface {
	ExchangeCode(ctx context.Context, idp *protocolShared.IdentityProviderConfig, code string) (IdpAuthentication, error)
}

type TokenDecoder interface {
	DecodeKid(token string) (string, error)
}

type ClientCredentialsSessionRepository interface {
	Save(ctx context.Context, sess *ClientCredentialsSession) error
	FindById(ctx context.Context, id ids.SessionId) (*ClientCredentialsSession, error)
}

type AccessTokenGenerator interface {
	GenerateAccessToken(ctx context.Context, session *BaseSession, identity *UserIdentityProjection) (string, error)
}

type IDTokenGenerator interface {
	GenerateIDToken(ctx context.Context, session *BaseSession, identity *UserIdentityProjection) (string, error)
}
