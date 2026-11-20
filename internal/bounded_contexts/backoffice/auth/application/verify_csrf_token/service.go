package verify_csrf_token

import (
	"context"
	"crypto/subtle"

	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

type SessionTokenReader interface {
	FindActiveCSRFToken(context.Context, ids.SessionId) (string, error)
}

type Service struct {
	sessions SessionTokenReader
}

func NewService(sessions SessionTokenReader) *Service {
	return &Service{sessions: sessions}
}

func (s *Service) Verify(ctx context.Context, id ids.SessionId, presented string) error {
	if presented == "" {
		return app_errors.ForbiddenError("INVALID_CSRF_TOKEN")
	}
	expected, err := s.sessions.FindActiveCSRFToken(ctx, id)
	if err != nil {
		return app_errors.InternalError("CSRF_TOKEN_VERIFICATION_FAILED").WithCause(err, false)
	}
	if subtle.ConstantTimeCompare([]byte(expected), []byte(presented)) != 1 {
		return app_errors.ForbiddenError("INVALID_CSRF_TOKEN")
	}
	return nil
}
