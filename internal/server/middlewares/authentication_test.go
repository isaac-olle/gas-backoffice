package middlewares_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	verifyJwt "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/key_pairs/application/verify_jwt"
	"github.com/isaac-olle/gas-backoffice/internal/server/middlewares"
	"github.com/isaac-olle/gas-backoffice/internal/server/middlewares/authenticators"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/query"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/authors"
)

const (
	frontendOrigin = "https://backoffice.example.com"
	userSub        = "b1b2c3d4-e5f6-7890-abcd-000000000001"
)

// recordingBus answers verify_jwt queries and records whether verification was
// reached, so tests can tell a token was extracted from the expected transport.
type recordingBus struct {
	result any
	err    error
	called bool
}

func (b *recordingBus) Dispatch(_ context.Context, q query.Query) (any, error) {
	if _, ok := q.(*verifyJwt.Query); ok {
		b.called = true
	}
	return b.result, b.err
}

func (b *recordingBus) Register(query.Type, query.Handler) {}

func validClaims() map[string]any {
	return map[string]any{"sub": userSub, "system_role": vo.SystemAdministratorClaim, "scope": "users:read"}
}

func TestAuthentication(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		origin     string
		bearer     string
		cookie     string
		busResult  any
		busErr     error
		wantNext   bool
		wantVerify bool
	}{
		{name: "cookie_origin_uses_cookie", origin: frontendOrigin, cookie: "cookie-token", busResult: validClaims(), wantNext: true, wantVerify: true},
		{name: "cookie_origin_ignores_bearer", origin: frontendOrigin, bearer: "bearer-token", wantNext: false},
		{name: "cookie_origin_without_cookie_aborts", origin: frontendOrigin, wantNext: false},
		{name: "other_origin_uses_bearer", origin: "https://evil.example.com", bearer: "bearer-token", busResult: validClaims(), wantNext: true, wantVerify: true},
		{name: "other_origin_ignores_cookie", origin: "https://evil.example.com", cookie: "cookie-token", wantNext: false},
		{name: "no_origin_uses_bearer", bearer: "bearer-token", busResult: validClaims(), wantNext: true, wantVerify: true},
		{name: "no_origin_without_bearer_aborts", wantNext: false},
		{name: "verification_error_aborts", bearer: "bearer-token", busErr: errors.New("invalid token"), wantNext: false, wantVerify: true},
		{name: "invalid_system_role_aborts", bearer: "bearer-token", busResult: map[string]any{"sub": userSub, "system_role": "nope"}, wantNext: false, wantVerify: true},
		{name: "missing_system_role_aborts", bearer: "bearer-token", busResult: map[string]any{"sub": userSub}, wantNext: false, wantVerify: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			gin.SetMode(gin.TestMode)
			bus := &recordingBus{result: tt.busResult, err: tt.busErr}

			nextCalled := false
			var gotAuthor *authors.UserAuthor
			r := gin.New()
			r.GET("/test", middlewares.Authentication(bus, "issuer", []string{frontendOrigin}), func(c *gin.Context) {
				nextCalled = true
				gotAuthor, _ = vo.AuthorFromContext[*authors.UserAuthor](c.Request.Context())
				c.Status(http.StatusOK)
			})

			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			if tt.origin != "" {
				req.Header.Set("Origin", tt.origin)
			}
			if tt.bearer != "" {
				req.Header.Set("Authorization", "Bearer "+tt.bearer)
			}
			if tt.cookie != "" {
				req.AddCookie(&http.Cookie{Name: authenticators.AccessTokenCookieName, Value: tt.cookie})
			}
			r.ServeHTTP(httptest.NewRecorder(), req)

			assert.Equal(t, tt.wantNext, nextCalled)
			assert.Equal(t, tt.wantVerify, bus.called)
			if tt.wantNext {
				require.NotNil(t, gotAuthor)
				assert.Equal(t, userSub, gotAuthor.Id.String())
				assert.Equal(t, userSub, gotAuthor.AdminId.String())
			}
		})
	}
}
