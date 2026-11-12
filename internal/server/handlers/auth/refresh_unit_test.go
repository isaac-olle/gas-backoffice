package auth_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/auth/application/refresh_session"
	"github.com/isaac-olle/gas-backoffice/internal/server/handlers/auth"
	"github.com/isaac-olle/gas-backoffice/internal/server/handlers/handlerstest"
	"github.com/isaac-olle/gas-backoffice/internal/server/middlewares/authenticators"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
)

func TestRefreshHandler_Execute(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		token   string
		busErr  error
		wantErr bool
	}{
		{name: "missing refresh token", wantErr: true},
		{name: "query error", token: "refresh-token", busErr: errors.New("query failed"), wantErr: true},
		{name: "success", token: "refresh-token"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			queries := handlerstest.NewQueryBus(refresh_session.Result{AccessToken: "access-token", RefreshToken: "new-refresh-token"}, tt.busErr)
			h := auth.NewRefreshHandler(queries)
			_, c := handlerstest.NewCtx(http.MethodPost, nil)
			if tt.token != "" {
				c.Request.AddCookie(&http.Cookie{Name: authenticators.RefreshTokenCookieName, Value: tt.token})
			}
			result, err := h.Execute(c)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.NoError(t, h.Success(c, result))
			assert.Equal(t, http.StatusNoContent, c.Writer.Status())
			assert.Contains(t, c.Writer.Header().Values("Set-Cookie"), authenticators.AccessTokenCookieName+"=access-token; Path=/; Max-Age=3600; HttpOnly; Secure; SameSite=None")
			assert.Contains(t, c.Writer.Header().Values("Set-Cookie"), authenticators.RefreshTokenCookieName+"=new-refresh-token; Path=/; Max-Age=2592000; HttpOnly; Secure; SameSite=None")
		})
	}
}

func TestRefreshHandlerMissingTokenCause(t *testing.T) {
	for _, emptyCookie := range []bool{false, true} {
		name, reason := "missing cookie", "missing"
		if emptyCookie {
			name, reason = "empty cookie", "empty"
		}
		t.Run(name, func(t *testing.T) {
			_, ctx := handlerstest.NewCtx(http.MethodPost, nil)
			if emptyCookie {
				ctx.Request.AddCookie(&http.Cookie{Name: authenticators.RefreshTokenCookieName, Value: ""})
			}
			_, err := auth.NewRefreshHandler(nil).Execute(ctx)
			var appError *app_errors.AppError
			require.ErrorAs(t, err, &appError)
			require.Equal(t, app_errors.ErrorCode("MISSING_REFRESH_TOKEN"), appError.Code)
			require.Equal(t, reason, appError.Details["reason"])
			require.Len(t, appError.PublicCauses(), 1)
			require.Contains(t, appError.PublicCauses()[0], "refresh_token")
		})
	}
}
