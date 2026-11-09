package refresh_session_test

import (
	"context"
	"errors"
	"testing"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/auth"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/auth/application/refresh_session"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"github.com/stretchr/testify/require"
)

type refreshVerifierStub struct {
	claims *auth.RefreshTokenClaims
	err    error
	token  string
}

func (v *refreshVerifierStub) Verify(_ context.Context, token string) (*auth.RefreshTokenClaims, error) {
	v.token = token
	return v.claims, v.err
}

type refreshAccessGenerator struct{ err error }

func (g *refreshAccessGenerator) Generate(context.Context, *auth.Session) (string, error) {
	return "access-token", g.err
}

type refreshTokenGenerator struct{ err error }

func (g *refreshTokenGenerator) GenerateRefreshToken(context.Context, *auth.Session) (string, error) {
	return "refresh-token", g.err
}

type refreshSessions struct {
	active bool
	err    error
	checks int
}

func (*refreshSessions) SaveSession(context.Context, *auth.Session) error { return nil }
func (r *refreshSessions) IsActive(context.Context, ids.SessionId) (bool, error) {
	r.checks++
	return r.active, r.err
}
func (*refreshSessions) Logout(context.Context, ids.SessionId) error { return nil }

type refreshBus struct{ events int }

func (b *refreshBus) Publish(_ context.Context, e []event.Event) { b.events += len(e) }

func TestRefresh(t *testing.T) {
	valid := &auth.RefreshTokenClaims{SessionID: ids.NewSessionId(), UserID: ids.NewUserId().String(), SystemRole: "system_administrator", Picture: "picture"}
	tests := []struct {
		name                                        string
		claims                                      *auth.RefreshTokenClaims
		verifyErr, activeErr, accessErr, refreshErr error
		active                                      bool
		checkSession                                bool
		wantErr                                     bool
	}{
		{name: "rotates tokens for an active session", claims: valid, active: true, checkSession: true},
		{name: "accepts valid claims without a session repository", claims: valid},
		{name: "rejects invalid refresh tokens", claims: valid, verifyErr: errors.New("invalid token"), wantErr: true},
		{name: "rejects missing user ids", claims: &auth.RefreshTokenClaims{SessionID: ids.NewSessionId(), SystemRole: "system_administrator"}, wantErr: true},
		{name: "rejects missing roles", claims: &auth.RefreshTokenClaims{SessionID: ids.NewSessionId(), UserID: ids.NewUserId().String()}, wantErr: true},
		{name: "rejects inactive sessions", claims: valid, active: false, checkSession: true, wantErr: true},
		{name: "rejects session repository errors", claims: valid, activeErr: errors.New("lookup failed"), checkSession: true, wantErr: true},
		{name: "rejects invalid role claims", claims: &auth.RefreshTokenClaims{SessionID: ids.NewSessionId(), UserID: ids.NewUserId().String(), SystemRole: "invalid"}, wantErr: true},
		{name: "wraps access token generation errors", claims: valid, active: true, checkSession: true, accessErr: errors.New("sign failed"), wantErr: true},
		{name: "wraps refresh token generation errors", claims: valid, active: true, checkSession: true, refreshErr: errors.New("sign failed"), wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			verifier := &refreshVerifierStub{claims: tt.claims, err: tt.verifyErr}
			access, refresh, bus := &refreshAccessGenerator{err: tt.accessErr}, &refreshTokenGenerator{err: tt.refreshErr}, &refreshBus{}
			var sessions auth.SessionRepository
			var sessionStub *refreshSessions
			if tt.checkSession {
				sessionStub = &refreshSessions{active: tt.active, err: tt.activeErr}
				sessions = sessionStub
			}
			result, err := refresh_session.NewService(verifier, access, refresh, sessions, bus).Refresh(context.Background(), refresh_session.NewQuery("refresh-token"))
			require.Equal(t, tt.wantErr, err != nil)
			require.Equal(t, "refresh-token", verifier.token)
			if sessionStub != nil {
				require.Equal(t, 1, sessionStub.checks)
			}
			if tt.wantErr {
				require.Zero(t, bus.events)
				return
			}
			got, ok := result.(refresh_session.Result)
			require.True(t, ok)
			require.Equal(t, refresh_session.Result{AccessToken: "access-token", RefreshToken: "refresh-token"}, got)
			require.Equal(t, 1, bus.events)
		})
	}
}
