package auth

import (
	"context"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/auth/vos"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

// SessionRepository persists the server-side state that makes JWT sessions
// revocable. Access and refresh tokens share the same session ID.
type SessionRepository interface {
	SaveSession(context.Context, *Session) error
	IsActive(context.Context, ids.SessionId) (bool, error)
	Logout(context.Context, ids.SessionId) error
}

type IdpConnector interface {
	Connect(ctx context.Context, idp *IdpConfig) (string, error)
}

type AuthorizationState struct {
	CodeVerifier     string
	IdpReferenceName string
}

type StateDecoder interface {
	DecodeState(state string) (AuthorizationState, error)
}

type IdpConfigRepository interface {
	FindByReferenceName(ctx context.Context, referenceName vos.IdpReferenceName) (*IdpConfig, error)
}

type TokenExchangeResult struct {
	Sub           string
	Email         string
	Picture       string
	EmailVerified bool
}

type TokenEndpointClient interface {
	ExchangeCode(ctx context.Context, code, codeVerifier, idpReferenceName string) (TokenExchangeResult, error)
}

type TokenGenerator interface {
	Generate(ctx context.Context, session *Session) (string, error)
}

type RefreshTokenGenerator interface {
	GenerateRefreshToken(ctx context.Context, session *Session) (string, error)
}

type RefreshTokenClaims struct {
	SessionID   ids.SessionId
	UserID      string
	SystemRole  string
	Picture     string
	Memberships []vo.SubscriptionMembership
}

type RefreshTokenVerifier interface {
	Verify(ctx context.Context, token string) (*RefreshTokenClaims, error)
}
