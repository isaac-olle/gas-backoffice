package authentication_methods_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas-backoffice/internal/server/handlers/handlerstest"
	"github.com/isaac-olle/gas-backoffice/internal/server/handlers/protocol/authentication_methods"
)

func bodyCredCtx(form url.Values, alreadyAuthenticated bool) *gin.Context {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodPost, "/", strings.NewReader(form.Encode()))
	c.Request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if alreadyAuthenticated {
		authentication_methods.MarkAuthenticated(c)
	}
	return c
}

func TestBodyCredentialsMethod_CanProcess(t *testing.T) {
	t.Parallel()

	full := url.Values{"client_id": {"abc"}, "client_secret": {"secret"}}

	tests := []struct {
		name          string
		form          url.Values
		contentType   string
		authenticated bool
		wantMatch     bool
		wantErr       bool
	}{
		{
			name:      "wrong_content_type",
			form:      full,
			wantMatch: false,
		},
		{
			name:      "missing_client_id",
			form:      url.Values{"client_secret": {"secret"}},
			wantMatch: false,
		},
		{
			name:      "missing_client_secret",
			form:      url.Values{"client_id": {"abc"}},
			wantMatch: false,
		},
		{
			name:      "all_present_not_authenticated",
			form:      full,
			wantMatch: true,
		},
		{
			name:          "all_present_already_authenticated",
			form:          full,
			authenticated: true,
			wantMatch:     false,
			wantErr:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			m := authentication_methods.NewBodyCredentialsMethod(handlerstest.NewQueryBus(nil, nil))

			var c *gin.Context
			if tt.name == "wrong_content_type" {
				w := httptest.NewRecorder()
				c, _ = gin.CreateTestContext(w)
				c.Request, _ = http.NewRequest(http.MethodPost, "/", strings.NewReader(tt.form.Encode()))
				c.Request.Header.Set("Content-Type", "application/json")
			} else {
				c = bodyCredCtx(tt.form, tt.authenticated)
			}

			match, err := m.CanProcess(c)
			assert.Equal(t, tt.wantMatch, match)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestBodyCredentialsMethod_Authenticate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		form       url.Values
		busErr     error
		wantErr    bool
		wantClient string
	}{
		{
			name:    "query_bus_error",
			form:    url.Values{"client_id": {"abc"}, "client_secret": {"secret"}},
			busErr:  assert.AnError,
			wantErr: true,
		},
		{
			name:       "success",
			form:       url.Values{"client_id": {"abc"}, "client_secret": {"secret"}},
			wantClient: "abc",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			m := authentication_methods.NewBodyCredentialsMethod(handlerstest.NewQueryBus(nil, tt.busErr))
			c := bodyCredCtx(tt.form, false)

			clientID, err := m.Authenticate(c)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.wantClient, clientID)
		})
	}
}
