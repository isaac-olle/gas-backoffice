package auth_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/auth"
	authHandlers "github.com/isaac-olle/gas-backoffice/internal/server/handlers/auth"
	"github.com/isaac-olle/gas-backoffice/internal/server/middlewares"
	"github.com/isaac-olle/gas-backoffice/internal/server/middlewares/authenticators"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

type csrfTokenGetter struct {
	sessionID ids.SessionId
}

func (g *csrfTokenGetter) Get(_ context.Context, id ids.SessionId) (string, error) {
	g.sessionID = id
	return "session-csrf", nil
}

type csrfRefreshVerifier struct {
	sessionID ids.SessionId
}

func (v csrfRefreshVerifier) Verify(_ context.Context, _ string) (*auth.RefreshTokenClaims, error) {
	return &auth.RefreshTokenClaims{SessionID: v.sessionID}, nil
}

func TestCSRFHandlerRecoversWithRefreshCookie(t *testing.T) {
	sessionID := ids.NewSessionId()
	tokens := &csrfTokenGetter{}
	bus := &stubQueryBus{err: errors.New("expired access token")}
	handler := authHandlers.NewCSRFHandler(tokens, csrfRefreshVerifier{sessionID}, bus, "gas")
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/v1/backoffice/auth/csrf", nil)
	c.Request.AddCookie(&http.Cookie{Name: authenticators.AccessTokenCookieName, Value: "expired"})
	c.Request.AddCookie(&http.Cookie{Name: authenticators.RefreshTokenCookieName, Value: "refresh"})

	result, err := handler.Execute(c)
	require.NoError(t, err)
	require.Equal(t, sessionID, tokens.sessionID)
	require.NoError(t, handler.Success(c, result))
	require.Equal(t, http.StatusOK, w.Code)
	require.Empty(t, w.Result().Cookies())
	require.Equal(t, "no-store", w.Header().Get("Cache-Control"))
	require.JSONEq(t, `{"csrf_token":"session-csrf"}`, w.Body.String())
}

func TestCSRFRequiresAllowedOrigin(t *testing.T) {
	const frontend = "https://backoffice.example.com"
	for _, origin := range []string{"", "null", "https://other.example.com", frontend} {
		t.Run(origin, func(t *testing.T) {
			reached := false
			router := gin.New()
			router.GET("/csrf", middlewares.RequireAllowedOrigin([]string{frontend}), func(c *gin.Context) {
				reached = true
				c.Status(http.StatusOK)
			})
			req := httptest.NewRequest(http.MethodGet, "/csrf", nil)
			if origin != "" {
				req.Header.Set("Origin", origin)
			}
			router.ServeHTTP(httptest.NewRecorder(), req)
			require.Equal(t, origin == frontend, reached)
		})
	}
}
