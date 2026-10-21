package auth

import (
	"context"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/auth/vos"
)

type IdpConnector interface {
	Connect(ctx context.Context, idp *IdpConfig) (string, error)
}

type StateDecoder interface {
	DecodeState(state string) (codeVerifier string, err error)
}

type IdpConfigRepository interface {
	FindByReferenceName(ctx context.Context, referenceName vos.IdpReferenceName) (*IdpConfig, error)
}

type TokenExchangeResult struct {
	Sub           string
	Email         string
	EmailVerified bool
}

type TokenEndpointClient interface {
	ExchangeCode(ctx context.Context, code, codeVerifier string) (TokenExchangeResult, error)
}

type TokenGenerator interface {
	Generate(ctx context.Context, session *Session) (string, error)
}
