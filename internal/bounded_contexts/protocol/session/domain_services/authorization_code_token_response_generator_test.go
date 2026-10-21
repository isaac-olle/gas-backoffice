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

	t.Run("includes an id_token when the session has the openid scope", func(t *testing.T) {
		t.Parallel()
		sess := makeSession(t) // scope: "openid", per mother_test.go
		gen := domain_services.NewAuthorizationCodeTokenResponseGenerator(
			&stubAccessTokenGenerator{token: "access-token"},
			&stubIDTokenGenerator{token: "id-token"},
		)

		resp, err := gen.Generate(context.Background(), sess, identity)

		require.NoError(t, err)
		assert.Equal(t, "access-token", resp.AccessToken)
		assert.Equal(t, "id-token", resp.IDToken)
	})

	t.Run("does not request an id_token when the session lacks the openid scope", func(t *testing.T) {
		t.Parallel()
		sess := makeSessionWithScope(t, "read")
		idTokenGen := &stubIDTokenGenerator{token: "id-token"}
		gen := domain_services.NewAuthorizationCodeTokenResponseGenerator(
			&stubAccessTokenGenerator{token: "access-token"},
			idTokenGen,
		)

		resp, err := gen.Generate(context.Background(), sess, identity)

		require.NoError(t, err)
		assert.Equal(t, "access-token", resp.AccessToken)
		assert.Empty(t, resp.IDToken)
	})

	t.Run("returns error when generating the access token fails", func(t *testing.T) {
		t.Parallel()
		sess := makeSession(t)
		gen := domain_services.NewAuthorizationCodeTokenResponseGenerator(
			&stubAccessTokenGenerator{err: assert.AnError},
			&stubIDTokenGenerator{},
		)

		_, err := gen.Generate(context.Background(), sess, identity)

		require.Error(t, err)
	})

	t.Run("returns error when generating the id_token fails", func(t *testing.T) {
		t.Parallel()
		sess := makeSession(t)
		gen := domain_services.NewAuthorizationCodeTokenResponseGenerator(
			&stubAccessTokenGenerator{token: "access-token"},
			&stubIDTokenGenerator{err: assert.AnError},
		)

		_, err := gen.Generate(context.Background(), sess, identity)

		require.Error(t, err)
	})
}
