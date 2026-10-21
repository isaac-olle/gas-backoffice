package session_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session"
	sessionvos "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session/vos"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
)

const (
	pr1 = "pr-1"
	pr2 = "pr-2"
)

func newFullClient() *session.Client {
	return session.NewClient(
		uuid.New().String(), "TestClient",
		uuid.New().String(),
		[]string{"authorization_code", "client_credentials"},
		[]string{"https://example.com/cb"},
		[]*session.Scope{
			session.NewScope(uuid.New().String(), "openid", pr1, "Resource A"),
			session.NewScope(uuid.New().String(), "profile", pr1, "Resource A"),
		},
		false,
	)
}

func TestClient_HasGrantType(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		grantType string
		want      bool
	}{
		{"returns true for registered grant type authorization_code", "authorization_code", true},
		{"returns true for registered grant type client_credentials", "client_credentials", true},
		{"returns false for unregistered grant type", "implicit", false},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			c := newFullClient()
			assert.Equal(t, tc.want, c.HasGrantType(tc.grantType))
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

func TestClient_ResolveProtectedResourceId(t *testing.T) {
	t.Parallel()

	clientWithScopes := func(scopes []*session.Scope) *session.Client {
		return session.NewClient(
			uuid.New().String(), "TC",
			uuid.New().String(),
			[]string{"client_credentials"},
			[]string{},
			scopes,
			false,
		)
	}

	tests := []struct {
		name            string
		scopes          []*session.Scope
		requestedScopes []string
		wantPrId        string
		wantErrCode     app_errors.ErrorCode
	}{
		{
			name: "resolves single scope to its PR",
			scopes: []*session.Scope{
				session.NewScope(uuid.New().String(), "users:read", pr1, "Resource A"),
			},
			requestedScopes: []string{"users:read"},
			wantPrId:        pr1,
		},
		{
			name: "resolves multiple scopes from same PR",
			scopes: []*session.Scope{
				session.NewScope(uuid.New().String(), "users:read", pr1, "Resource A"),
				session.NewScope(uuid.New().String(), "users:write", pr1, "Resource A"),
			},
			requestedScopes: []string{"users:read", "users:write"},
			wantPrId:        pr1,
		},
		{
			name: "rejects scope name present in multiple PRs",
			scopes: []*session.Scope{
				session.NewScope(uuid.New().String(), "users:read", pr1, "Resource A"),
				session.NewScope(uuid.New().String(), "users:read", pr2, "Resource B"),
			},
			requestedScopes: []string{"users:read"},
			wantErrCode:     session.ErrAmbiguousScope,
		},
		{
			name: "rejects scopes from different PRs",
			scopes: []*session.Scope{
				session.NewScope(uuid.New().String(), "users:read", pr1, "Resource A"),
				session.NewScope(uuid.New().String(), "orders:read", pr2, "Resource B"),
			},
			requestedScopes: []string{"users:read", "orders:read"},
			wantErrCode:     session.ErrScopesFromMultipleResources,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			c := clientWithScopes(tc.scopes)
			prId, err := c.ResolveProtectedResourceId(tc.requestedScopes)

			if tc.wantErrCode != "" {
				var appErr *app_errors.AppError
				require.ErrorAs(t, err, &appErr)
				assert.Equal(t, tc.wantErrCode, appErr.Code)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.wantPrId, prId)
		})
	}
}

func TestClient_CheckClientCredentials(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		grantTypes   []string
		scopes       []*session.Scope
		requestScope string
		wantPrId     string
		wantErrCode  app_errors.ErrorCode
	}{
		{
			name:       "returns PR id for valid single scope",
			grantTypes: []string{"client_credentials"},
			scopes: []*session.Scope{
				session.NewScope(uuid.New().String(), "openid", pr1, "Resource A"),
			},
			requestScope: "openid",
			wantPrId:     pr1,
		},
		{
			name:       "returns PR id for multiple scopes from same PR",
			grantTypes: []string{"client_credentials"},
			scopes: []*session.Scope{
				session.NewScope(uuid.New().String(), "openid", pr1, "Resource A"),
				session.NewScope(uuid.New().String(), "profile", pr1, "Resource A"),
			},
			requestScope: "openid profile",
			wantPrId:     pr1,
		},
		{
			name:       "rejects client without client_credentials grant type",
			grantTypes: []string{"authorization_code"},
			scopes: []*session.Scope{
				session.NewScope(uuid.New().String(), "openid", pr1, "Resource A"),
			},
			requestScope: "openid",
			wantErrCode:  session.ErrGrantTypeNotAllowed,
		},
		{
			name:       "rejects scope not registered for client",
			grantTypes: []string{"client_credentials"},
			scopes: []*session.Scope{
				session.NewScope(uuid.New().String(), "openid", pr1, "Resource A"),
			},
			requestScope: "email",
			wantErrCode:  session.ErrInvalidScope,
		},
		{
			name:       "rejects ambiguous scope name across PRs",
			grantTypes: []string{"client_credentials"},
			scopes: []*session.Scope{
				session.NewScope(uuid.New().String(), "users:read", pr1, "Resource A"),
				session.NewScope(uuid.New().String(), "users:read", pr2, "Resource B"),
			},
			requestScope: "users:read",
			wantErrCode:  session.ErrAmbiguousScope,
		},
		{
			name:       "rejects scopes from different PRs",
			grantTypes: []string{"client_credentials"},
			scopes: []*session.Scope{
				session.NewScope(uuid.New().String(), "users:read", pr1, "Resource A"),
				session.NewScope(uuid.New().String(), "orders:read", pr2, "Resource B"),
			},
			requestScope: "users:read orders:read",
			wantErrCode:  session.ErrScopesFromMultipleResources,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			c := session.NewClient(
				uuid.New().String(), "TC",
				uuid.New().String(),
				tc.grantTypes, []string{},
				tc.scopes,
				false,
			)

			prId, err := c.CheckClientCredentials(tc.requestScope)
			if tc.wantErrCode != "" {
				var appErr *app_errors.AppError
				require.ErrorAs(t, err, &appErr)
				assert.Equal(t, tc.wantErrCode, appErr.Code)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.wantPrId, prId)
		})
	}
}

func TestClientCredentialsGrantTypeConstant(t *testing.T) {
	t.Parallel()
	assert.Equal(t, "client_credentials", sessionvos.ClientCredentialsGrantType.String())
}

func TestNewScope(t *testing.T) {
	t.Parallel()
	s := session.NewScope("scope-1", "read:users", "pr-1", "My Resource")
	assert.Equal(t, "scope-1", s.ScopeID)
	assert.Equal(t, "read:users", s.ScopeName)
	assert.Equal(t, "pr-1", s.ProtectedResourceID)
	assert.Equal(t, "My Resource", s.ProtectedResourceName)
}
