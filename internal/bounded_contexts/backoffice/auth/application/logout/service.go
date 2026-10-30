package logout

import (
	"context"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/auth"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
)

type Service struct{ sessions auth.SessionRepository }

func NewService(sessions auth.SessionRepository) *Service { return &Service{sessions: sessions} }

func (s *Service) Logout(ctx context.Context, cmd *Command) error {
	if err := s.sessions.Logout(ctx, cmd.sessionID); err != nil {
		return app_errors.InternalError(auth.ErrCannotLogoutSession).WithCause(err, false)
	}
	return nil
}
