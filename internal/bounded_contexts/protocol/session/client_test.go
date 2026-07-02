package session_test

import (
	"github.com/isaac-olle/gas/internal/shared/app_errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas/internal/bounded_contexts/protocol/session"
	sessionvos "github.com/isaac-olle/gas/internal/bounded_contexts/protocol/session/vos"
)

func newFullClient() *session.Client {
	return session.NewClient(
		uuid.New().String(), "TestClient",
		uuid.New().String(), uuid.New().String(),
		[]string{"authorization_code", "client_credentials"},
		[]string{"https://example.com/cb"},
		[]*session.Scope{
			{ScopeName: "openid"},
			{ScopeName: "profile"},
		},
	)
}

func TestClient_HasFlow(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		flow string
		want bool
	}{
		{"returns true for registered flow authorization_code", "authorization_code", true},
		{"returns true for registered flow client_credentials", "client_credentials", true},
		{"returns false for unregistered flow", "pkce", false},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			c := newFullClient()
			assert.Equal(t, tc.want, c.HasFlow(tc.flow))
		})
	}
}

func TestClient_HasRedirectURI(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		uri  string
		want bool
	}{
		{"returns true for registered uri", "https://example.com/cb", true},
		{"returns false for unregistered uri", "https://evil.com/cb", false},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			c := newFullClient()
			assert.Equal(t, tc.want, c.HasRedirectURI(tc.uri))
		})
	}
}

func TestClient_HasScope(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		scope string
		want  bool
	}{
		{"returns true for registered scope openid", "openid", true},
		{"returns true for registered scope profile", "profile", true},
		{"returns false for unregistered scope email", "email", false},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			c := newFullClient()
			assert.Equal(t, tc.want, c.HasScope(tc.scope))
		})
	}
}

func TestClient_CheckClientCredentials(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		flows        []string
		scopes       []string
		requestScope string
		wantErr      bool
		wantErrCode  app_errors.ErrorCode
	}{
		{
			name:         "accepts valid flow and scopes",
			flows:        []string{"client_credentials"},
			scopes:       []string{"openid", "profile"},
			requestScope: "openid profile",
		},
		{
			name:         "rejects client without client_credentials flow",
			flows:        []string{"authorization_code"},
			scopes:       []string{"openid"},
			requestScope: "openid",
			wantErr:      true,
			wantErrCode:  session.ErrFlowNotAllowed,
		},
		{
			name:         "rejects scope not registered for client",
			flows:        []string{"client_credentials"},
			scopes:       []string{"openid"},
			requestScope: "email",
			wantErr:      true,
			wantErrCode:  session.ErrInvalidScope,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			scopeObjs := make([]*session.Scope, len(tc.scopes))
			for i, s := range tc.scopes {
				scopeObjs[i] = &session.Scope{ScopeName: s}
			}
			c := session.NewClient(
				uuid.New().String(), "TC",
				uuid.New().String(), uuid.New().String(),
				tc.flows, []string{},
				scopeObjs,
			)

			err := c.CheckClientCredentials(tc.requestScope)
			if tc.wantErr {
				var appErr *app_errors.AppError
				require.ErrorAs(t, err, &appErr)
				assert.Equal(t, tc.wantErrCode, appErr.Code)
				return
			}
			assert.NoError(t, err)
		})
	}
}

func TestClientCredentialsFlowConstant(t *testing.T) {
	t.Parallel()
	assert.Equal(t, "client_credentials", sessionvos.ClientCredentialsFlow)
}

func TestNewScope(t *testing.T) {
	t.Parallel()
	s := session.NewScope("scope-1", "read:users", "pr-1", "My Resource")
	assert.Equal(t, "scope-1", s.ScopeID)
	assert.Equal(t, "read:users", s.ScopeName)
	assert.Equal(t, "pr-1", s.ProtectedResourceID)
	assert.Equal(t, "My Resource", s.ProtectedResourceName)
}
