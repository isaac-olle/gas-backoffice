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
		uuid.New().String(),
		"read",
		uuid.New().String(),
	)
	require.NoError(t, err)
	return sess
}

func TestClientCredentialsTokenResponseGenerator_Generate(t *testing.T) {
	t.Parallel()

	identity := session.NewUserIdentityProjection("client-subject", "")

	t.Run("returns only an access token, never an id_token", func(t *testing.T) {
		t.Parallel()
		sess := makeClientCredentialsSession(t)
		gen := domain_services.NewClientCredentialsTokenResponseGenerator(
			&stubAccessTokenGenerator{token: "access-token"},
		)

		resp, err := gen.Generate(context.Background(), sess, identity)

		require.NoError(t, err)
		assert.Equal(t, "access-token", resp.AccessToken)
		assert.Empty(t, resp.IDToken)
	})

	t.Run("returns error when generating the access token fails", func(t *testing.T) {
		t.Parallel()
		sess := makeClientCredentialsSession(t)
		gen := domain_services.NewClientCredentialsTokenResponseGenerator(
			&stubAccessTokenGenerator{err: assert.AnError},
		)

		_, err := gen.Generate(context.Background(), sess, identity)

		require.Error(t, err)
	})
}
