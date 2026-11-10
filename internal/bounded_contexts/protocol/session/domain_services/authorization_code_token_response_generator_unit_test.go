package domain_services_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session/domain_services"
)

type stubAccessTokenGenerator struct {
	token string
	err   error
}

func (s *stubAccessTokenGenerator) GenerateAccessToken(_ context.Context, _ *session.BaseSession, _ *session.UserIdentityProjection) (string, error) {
	return s.token, s.err
}

type stubIDTokenGenerator struct {
	token string
	err   error
}

func (s *stubIDTokenGenerator) GenerateIDToken(_ context.Context, _ *session.BaseSession, _ *session.UserIdentityProjection) (string, error) {
	return s.token, s.err
}

func TestAuthorizationCodeTokenResponseGenerator_Generate(t *testing.T) {
	t.Parallel()

	identity := session.NewUserIdentityProjection("user-subject", "")

	tests := []struct {
		name          string
		makeSession   func(t *testing.T) *session.AuthorizationCodeSession
		accessToken   string
		accessErr     error
		idToken       string
		idErr         error
		wantErr       bool
		wantIDToken   string
		wantNoIDToken bool
	}{
		{
			name:        "includes an id_token when the session has the openid scope",
			makeSession: makeSession, // scope: "openid", per mother_test.go
			accessToken: "access-token",
			idToken:     "id-token",
			wantIDToken: "id-token",
		},
		{
			name:          "does not request an id_token when the session lacks the openid scope",
			makeSession:   func(t *testing.T) *session.AuthorizationCodeSession { return makeSessionWithScope(t, "read") },
			accessToken:   "access-token",
			idToken:       "id-token",
			wantNoIDToken: true,
		},
		{
			name:        "returns error when generating the access token fails",
			makeSession: makeSession,
			accessErr:   assert.AnError,
			wantErr:     true,
		},
		{
			name:        "returns error when generating the id_token fails",
			makeSession: makeSession,
			accessToken: "access-token",
			idErr:       assert.AnError,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			sess := tt.makeSession(t)
			gen := domain_services.NewAuthorizationCodeTokenResponseGenerator(
				&stubAccessTokenGenerator{token: tt.accessToken, err: tt.accessErr},
				&stubIDTokenGenerator{token: tt.idToken, err: tt.idErr},
			)

			resp, err := gen.Generate(context.Background(), sess, identity)

			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.accessToken, resp.AccessToken)
			if tt.wantNoIDToken {
				assert.Empty(t, resp.IDToken)
				return
			}
			assert.Equal(t, tt.wantIDToken, resp.IDToken)
		})
	}
}
