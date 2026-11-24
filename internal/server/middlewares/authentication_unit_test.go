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
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
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
	return map[string]any{"sub": userSub, "sid": userSub, "system_role": vo.SystemAdministratorClaim, "scope": "users:read"}
}

type csrfVerifier struct {
	sessionID string
	token     string
	err       error
}

func (v *csrfVerifier) Verify(_ context.Context, id ids.SessionId, token string) error {
	v.sessionID, v.token = id.String(), token
	return v.err
}

func TestAuthentication(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		origin     string
		bearer     string
		authHeader *string
		cookie     string
		busResult  any
		busErr     error
		wantNext   bool
		wantVerify bool
		method     string
		csrf       string
		csrfErr    error
		wantCSRF   bool
	}{
		{name: "cookie_origin_uses_cookie", origin: frontendOrigin, cookie: "cookie-token", busResult: validClaims(), wantNext: true, wantVerify: true},
		{name: "explicit_bearer_skips_csrf_even_from_frontend", method: http.MethodPost, origin: frontendOrigin, bearer: "bearer-token", busResult: validClaims(), wantNext: true, wantVerify: true},
		{name: "empty_authorization_does_not_fall_back_to_cookie", origin: frontendOrigin, authHeader: new(string), cookie: "cookie-token", busResult: validClaims()},
		{name: "cookie_write_checks_csrf", method: http.MethodPost, origin: frontendOrigin, cookie: "cookie-token", csrf: "csrf-token", busResult: validClaims(), wantNext: true, wantVerify: true, wantCSRF: true},
		{name: "csrf_failure_aborts_cookie_write", method: http.MethodPost, origin: frontendOrigin, cookie: "cookie-token", busResult: validClaims(), csrfErr: errors.New("missing CSRF token"), wantVerify: true, wantCSRF: true},
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
			csrf := &csrfVerifier{err: tt.csrfErr}
			cookie := authenticators.NewCookieAuthenticator(csrf)
			method := tt.method
			if method == "" {
				method = http.MethodGet
			}

			nextCalled := false
			var gotAuthor *vo.UserAuthor
			r := gin.New()
			r.Handle(method, "/test", middlewares.Authentication(bus, "issuer", []string{frontendOrigin}, cookie), func(c *gin.Context) {
				nextCalled = true
				gotAuthor, _ = vo.AuthorFromContext(c.Request.Context())
				c.Status(http.StatusOK)
			})

			req := httptest.NewRequest(method, "/test", nil)
			req.Header.Set("X-CSRF-Token", tt.csrf)
			if tt.origin != "" {
				req.Header.Set("Origin", tt.origin)
			}
			if tt.bearer != "" {
				req.Header.Set("Authorization", "Bearer "+tt.bearer)
			}
			if tt.authHeader != nil {
				req.Header.Set("Authorization", *tt.authHeader)
			}
			if tt.cookie != "" {
				req.AddCookie(&http.Cookie{Name: authenticators.AccessTokenCookieName, Value: tt.cookie})
			}
			r.ServeHTTP(httptest.NewRecorder(), req)

			assert.Equal(t, tt.wantNext, nextCalled)
			assert.Equal(t, tt.wantVerify, bus.called)
			if tt.wantCSRF {
				assert.Equal(t, userSub, csrf.sessionID)
				assert.Equal(t, tt.csrf, csrf.token)
			} else {
				assert.Empty(t, csrf.sessionID)
			}
			if tt.wantNext {
				require.NotNil(t, gotAuthor)
				assert.Equal(t, userSub, gotAuthor.Id.String())
				assert.Equal(t, userSub, gotAuthor.AdminId.String())
			}
		})
	}
}
