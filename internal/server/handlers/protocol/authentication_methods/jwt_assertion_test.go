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

	verifyAssertion "github.com/isaac-olle/gas/internal/bounded_contexts/security/client_assertion_keys/application/verify_assertion"
	"github.com/isaac-olle/gas/internal/server/handlers/handlerstest"
	"github.com/isaac-olle/gas/internal/server/handlers/protocol/authentication_methods"
)

const jwtBearerType = "urn:ietf:params:oauth:client-assertion-type:jwt-bearer"

func jwtCtx(form url.Values, alreadyAuthenticated bool) *gin.Context {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodPost, "/", strings.NewReader(form.Encode()))
	c.Request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if alreadyAuthenticated {
		authentication_methods.MarkAuthenticated(c)
	}
	return c
}

func TestJWTAssertionMethod_CanProcess(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		form          url.Values
		authenticated bool
		wantMatch     bool
		wantErr       bool
	}{
		{
			name:      "no_assertion_type",
			form:      url.Values{},
			wantMatch: false,
		},
		{
			name:      "wrong_assertion_type",
			form:      url.Values{"client_assertion_type": {"urn:other"}},
			wantMatch: false,
		},
		{
			name:      "correct_assertion_type_not_authenticated",
			form:      url.Values{"client_assertion_type": {jwtBearerType}},
			wantMatch: true,
		},
		{
			name:          "correct_assertion_type_already_authenticated",
			form:          url.Values{"client_assertion_type": {jwtBearerType}},
			authenticated: true,
			wantMatch:     false,
			wantErr:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			m := authentication_methods.NewJWTAssertionMethod(handlerstest.NewQueryBus(nil, nil))
			c := jwtCtx(tt.form, tt.authenticated)

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

func TestJWTAssertionMethod_Authenticate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		form       url.Values
		busResult  any
		busErr     error
		wantErr    bool
		wantClient string
	}{
		{
			name:    "missing_client_assertion",
			form:    url.Values{"client_assertion_type": {jwtBearerType}},
			wantErr: true,
		},
		{
			name:    "query_bus_error",
			form:    url.Values{"client_assertion_type": {jwtBearerType}, "client_assertion": {"the.jwt.token"}},
			busErr:  assert.AnError,
			wantErr: true,
		},
		{
			name:      "bus_returns_invalid_type",
			form:      url.Values{"client_assertion_type": {jwtBearerType}, "client_assertion": {"the.jwt.token"}},
			busResult: "not-a-claims-struct",
			wantErr:   true,
		},
		{
			name:      "bus_returns_empty_client_id",
			form:      url.Values{"client_assertion_type": {jwtBearerType}, "client_assertion": {"the.jwt.token"}},
			busResult: &verifyAssertion.ClientClaims{ClientID: ""},
			wantErr:   true,
		},
		{
			name:       "success",
			form:       url.Values{"client_assertion_type": {jwtBearerType}, "client_assertion": {"the.jwt.token"}},
			busResult:  &verifyAssertion.ClientClaims{ClientID: "client-xyz"},
			wantClient: "client-xyz",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			m := authentication_methods.NewJWTAssertionMethod(handlerstest.NewQueryBus(tt.busResult, tt.busErr))
			c := jwtCtx(tt.form, false)

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
