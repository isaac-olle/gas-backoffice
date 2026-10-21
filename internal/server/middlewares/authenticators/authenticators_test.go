package authenticators_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas-backoffice/internal/server/middlewares/authenticators"
)

func newCtx(t *testing.T, setup func(r *http.Request)) *gin.Context {
	t.Helper()
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	setup(c.Request)
	return c
}

func TestBearerAuthenticator(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		header  string
		want    string
		wantErr bool
	}{
		{name: "missing_header_errors", header: "", wantErr: true},
		{name: "non_bearer_header_errors", header: "Basic abc", wantErr: true},
		{name: "empty_bearer_errors", header: "Bearer ", wantErr: true},
		{name: "bearer_returns_token", header: "Bearer abc", want: "abc"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			c := newCtx(t, func(r *http.Request) {
				if tt.header != "" {
					r.Header.Set("Authorization", tt.header)
				}
			})
			token, err := authenticators.NewBearerAuthenticator().Token(c)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, token)
		})
	}
}

func TestCookieAuthenticator(t *testing.T) {
	t.Parallel()

	t.Run("missing_cookie_errors", func(t *testing.T) {
		t.Parallel()
		c := newCtx(t, func(r *http.Request) {})
		_, err := authenticators.NewCookieAuthenticator().Token(c)
		assert.Error(t, err)
	})

	t.Run("ignores_bearer_header", func(t *testing.T) {
		t.Parallel()
		c := newCtx(t, func(r *http.Request) { r.Header.Set("Authorization", "Bearer abc") })
		_, err := authenticators.NewCookieAuthenticator().Token(c)
		assert.Error(t, err)
	})

	t.Run("cookie_returns_token", func(t *testing.T) {
		t.Parallel()
		c := newCtx(t, func(r *http.Request) {
			r.AddCookie(&http.Cookie{Name: authenticators.AccessTokenCookieName, Value: "abc"})
		})
		token, err := authenticators.NewCookieAuthenticator().Token(c)
		require.NoError(t, err)
		assert.Equal(t, "abc", token)
	})
}
