package authentication_methods_test

import (
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas/internal/server/handlers/handlerstest"
	"github.com/isaac-olle/gas/internal/server/handlers/protocol/authentication_methods"
)

func basicAuthCtx(authHeader string, alreadyAuthenticated bool) *gin.Context {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodPost, "/", nil)
	if authHeader != "" {
		c.Request.Header.Set("Authorization", authHeader)
	}
	if alreadyAuthenticated {
		authentication_methods.MarkAuthenticated(c)
	}
	return c
}

func basicHeader(clientID, secret string) string {
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(clientID+":"+secret))
}

func TestBasicAuthMethod_CanProcess(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		authHeader    string
		authenticated bool
		wantMatch     bool
		wantErr       bool
	}{
		{
			name:      "no_authorization_header",
			wantMatch: false,
		},
		{
			name:       "bearer_token_header",
			authHeader: "Bearer sometoken",
			wantMatch:  false,
		},
		{
			name:       "basic_header_not_authenticated",
			authHeader: basicHeader("abc", "secret"),
			wantMatch:  true,
		},
		{
			name:          "basic_header_already_authenticated",
			authHeader:    basicHeader("abc", "secret"),
			authenticated: true,
			wantMatch:     false,
			wantErr:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			m := authentication_methods.NewBasicAuthMethod(handlerstest.NewQueryBus(nil, nil))
			c := basicAuthCtx(tt.authHeader, tt.authenticated)

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

func TestBasicAuthMethod_Authenticate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		authHeader string
		busErr     error
		wantErr    bool
		wantClient string
	}{
		{
			name:       "invalid_base64",
			authHeader: "Basic not-valid-base64!!!",
			wantErr:    true,
		},
		{
			name:       "no_colon_in_credentials",
			authHeader: "Basic " + base64.StdEncoding.EncodeToString([]byte("justid")),
			wantErr:    true,
		},
		{
			name:       "query_bus_error",
			authHeader: basicHeader("abc", "secret"),
			busErr:     assert.AnError,
			wantErr:    true,
		},
		{
			name:       "success",
			authHeader: basicHeader("abc", "secret"),
			wantClient: "abc",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			m := authentication_methods.NewBasicAuthMethod(handlerstest.NewQueryBus(nil, tt.busErr))
			c := basicAuthCtx(tt.authHeader, false)

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
