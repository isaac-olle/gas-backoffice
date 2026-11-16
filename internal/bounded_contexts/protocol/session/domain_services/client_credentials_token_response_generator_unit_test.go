package domain_services_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session/domain_services"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

func makeClientCredentialsSession(t *testing.T) *session.ClientCredentialsSession {
	t.Helper()
	sess, err := session.NewClientCredentialsSession(
		ids.NewSessionId(),
		uuid.New().String(),
		"read",
		uuid.New().String(),
	)
	require.NoError(t, err)
	return sess
}

func TestClientCredentialsTokenResponseGenerator_Generate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		accessToken string
		accessErr   error
		wantErr     bool
	}{
		{
			name:        "returns only an access token, never an id_token",
			accessToken: "access-token",
		},
		{
			name:      "returns error when generating the access token fails",
			accessErr: assert.AnError,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			sess := makeClientCredentialsSession(t)
			accessGenerator := &stubAccessTokenGenerator{token: tt.accessToken, err: tt.accessErr}
			gen := domain_services.NewClientCredentialsTokenResponseGenerator(
				accessGenerator,
			)

			resp, err := gen.Generate(context.Background(), sess)
			assert.Equal(t, sess.ClientID().String(), accessGenerator.claims.Subject)
			assert.Empty(t, accessGenerator.claims.IdpID)

			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.accessToken, resp.AccessToken)
			assert.Empty(t, resp.IDToken)
		})
	}
}
