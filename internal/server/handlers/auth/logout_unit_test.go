package auth_test

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas-backoffice/internal/server/handlers/auth"
	"github.com/isaac-olle/gas-backoffice/internal/server/handlers/handlerstest"
	"github.com/isaac-olle/gas-backoffice/internal/server/middlewares/authenticators"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

type logoutCSRFVerifier struct {
	sessionID ids.SessionId
	token     string
	err       error
}

func (v *logoutCSRFVerifier) Verify(_ context.Context, id ids.SessionId, token string) error {
	v.sessionID, v.token = id, token
	return v.err
}

func TestLogoutHandler_Execute(t *testing.T) {
	t.Parallel()
	sessionID := ids.NewSessionId().String()
	tests := []struct {
		name          string
		accessToken   string
		refreshToken  string
		invalidClaims bool
		wantQueries   int
		wantCommands  int
		csrfErr       error
	}{
		{name: "no token", wantQueries: 0, wantCommands: 0},
		{name: "access token", accessToken: "access-token", wantQueries: 1, wantCommands: 1},
		{name: "refresh token", refreshToken: "refresh-token", wantQueries: 1, wantCommands: 1},
		{name: "invalid claims", accessToken: "access-token", invalidClaims: true, wantQueries: 1, wantCommands: 0},
		{name: "csrf error blocks logout", refreshToken: "refresh-token", wantQueries: 1, csrfErr: errors.New("invalid CSRF token")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			claims := map[string]any{"sid": sessionID}
			var queryErr error
			if tt.invalidClaims {
				claims = map[string]any{"sid": 42}
			}
			queries := handlerstest.NewQueryBus(claims, queryErr)
			commands := handlerstest.NewCommandExecutor(nil)
			csrfVerifier := &logoutCSRFVerifier{err: tt.csrfErr}
			h := auth.NewLogoutHandler(commands, queries, "https://issuer.example", authenticators.NewCookieAuthenticator(csrfVerifier))
			_, c := handlerstest.NewCtx(http.MethodPost, nil)
			c.Request.Header.Set("X-CSRF-Token", "csrf-token")
			if tt.accessToken != "" {
				c.Request.AddCookie(&http.Cookie{Name: authenticators.AccessTokenCookieName, Value: tt.accessToken})
			}
			if tt.refreshToken != "" {
				c.Request.AddCookie(&http.Cookie{Name: authenticators.RefreshTokenCookieName, Value: tt.refreshToken})
			}

			_, err := h.Execute(c)
			queries.AssertNumberOfCalls(t, "Dispatch", tt.wantQueries)
			commands.AssertNumberOfCalls(t, "SingleExecution", tt.wantCommands)
			if tt.wantCommands > 0 || tt.csrfErr != nil {
				assert.Equal(t, sessionID, csrfVerifier.sessionID.String())
				assert.Equal(t, "csrf-token", csrfVerifier.token)
			}
			if tt.csrfErr != nil {
				require.ErrorIs(t, err, tt.csrfErr)
				return
			}
			require.NoError(t, err)
			require.NoError(t, h.Success(c, nil))
			assert.Equal(t, http.StatusNoContent, c.Writer.Status())
			assert.Contains(t, c.Writer.Header().Get("Set-Cookie"), authenticators.AccessTokenCookieName+"=")
		})
	}
}
