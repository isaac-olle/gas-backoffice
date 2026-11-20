package get_csrf_token

import (
	"context"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/auth"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

type SessionTokens interface {
	FindActiveCSRFToken(context.Context, ids.SessionId) (string, error)
	InitializeCSRFToken(context.Context, ids.SessionId, string) (string, error)
}

type Service struct {
	sessions  SessionTokens
	generator auth.CSRFTokenGenerator
}

func NewService(sessions SessionTokens, generator auth.CSRFTokenGenerator) *Service {
	return &Service{sessions: sessions, generator: generator}
}

func (s *Service) Get(ctx context.Context, id ids.SessionId) (string, error) {
	token, err := s.sessions.FindActiveCSRFToken(ctx, id)
	if err == nil && token == "" {
		token, err = s.sessions.InitializeCSRFToken(ctx, id, s.generator.GenerateCSRFToken())
	}
	if err != nil {
		return "", app_errors.InternalError("CSRF_TOKEN_LOOKUP_FAILED").WithCause(err, false)
	}
	if token == "" {
		return "", app_errors.UnauthorizedError("SESSION_REVOKED")
	}
	return token, nil
}
