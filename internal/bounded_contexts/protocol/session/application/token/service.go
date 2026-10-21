package token

import (
	"context"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
)

// Request holds the raw /oidc/token parameters, across every grant type.
// Each GrantTypeHandler reads only the fields it needs.
type Request struct {
	ClientID     string
	IsPublic     bool
	Code         string
	RedirectURI  string
	CodeVerifier string
	Scope        string
}

// GrantTypeHandler resolves and executes a single OAuth2 grant type.
type GrantTypeHandler interface {
	Supports(grantType string) bool
	Execute(ctx context.Context, req Request) (*session.TokenResponse, error)
}

// Service is the atomic use case for the token step: it picks the
// GrantTypeHandler that supports the requested grant_type and delegates to
// it, without going through the command/query bus. Adding a new grant type
// means adding a new GrantTypeHandler, not touching this service or the HTTP
// handler.
type Service struct {
	handlers []GrantTypeHandler
}

func NewService(handlers ...GrantTypeHandler) *Service {
	return &Service{handlers: handlers}
}

func (s *Service) Execute(ctx context.Context, grantType string, req Request) (*session.TokenResponse, error) {
	for _, h := range s.handlers {
		if h.Supports(grantType) {
			return h.Execute(ctx, req)
		}
	}
	return nil, app_errors.BadRequestError(session.ErrUnsupportedGrantType).
		WithDescription("unsupported grant_type: " + grantType)
}
