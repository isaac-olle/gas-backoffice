package logout

import (
	"context"
	"errors"
	"testing"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/auth"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"github.com/stretchr/testify/require"
)

type logoutSessions struct {
	got ids.SessionId
	err error
}

func (s *logoutSessions) SaveSession(context.Context, *auth.Session) error      { return nil }
func (s *logoutSessions) IsActive(context.Context, ids.SessionId) (bool, error) { return false, nil }
func (s *logoutSessions) Logout(_ context.Context, id ids.SessionId) error      { s.got = id; return s.err }

func TestLogout(t *testing.T) {
	sessionID := ids.NewSessionId()
	tests := []struct {
		name    string
		repoErr error
		wantErr bool
	}{
		{name: "logs out the requested session"},
		{name: "wraps repository failures", repoErr: errors.New("write failed"), wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &logoutSessions{err: tt.repoErr}
			err := NewService(repo).Logout(context.Background(), NewCommand(sessionID))
			require.Equal(t, tt.wantErr, err != nil)
			require.Equal(t, sessionID, repo.got)
		})
	}
}
