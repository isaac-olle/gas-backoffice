package vos_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	sessionvos "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session/vos"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

func validUserIdentityArgs() (string, string, string, sessionvos.IdpClaims, time.Time) {
	return ids.NewIdentityProviderId().String(),
		"id-token",
		"access-token",
		sessionvos.IdpClaims{Sub: "google-sub", Email: "user@example.com"},
		time.Now().Add(time.Hour)
}

func TestNewUserIdentity(t *testing.T) {
	t.Parallel()

	t.Run("builds a valid identity", func(t *testing.T) {
		t.Parallel()
		idpId, idToken, accessToken, claims, expiresAt := validUserIdentityArgs()
		refreshToken := "refresh-token"

		identity, err := sessionvos.NewUserIdentity(idpId, idToken, accessToken, &refreshToken, claims, expiresAt)

		require.NoError(t, err)
		assert.Equal(t, idpId, identity.IdpId().String())
		assert.Equal(t, idToken, identity.IDToken())
		assert.Equal(t, accessToken, identity.AccessToken())
		require.NotNil(t, identity.RefreshToken())
		assert.Equal(t, refreshToken, *identity.RefreshToken())
		assert.Equal(t, claims, identity.Claims())
		assert.WithinDuration(t, expiresAt, identity.ExpiresAt(), 0)
	})

	t.Run("allows a nil refresh token", func(t *testing.T) {
		t.Parallel()
		idpId, idToken, accessToken, claims, expiresAt := validUserIdentityArgs()

		identity, err := sessionvos.NewUserIdentity(idpId, idToken, accessToken, nil, claims, expiresAt)

		require.NoError(t, err)
		assert.Nil(t, identity.RefreshToken())
	})

	t.Run("rejects an invalid idp id", func(t *testing.T) {
		t.Parallel()
		_, idToken, accessToken, claims, expiresAt := validUserIdentityArgs()

		_, err := sessionvos.NewUserIdentity("not-a-uuid", idToken, accessToken, nil, claims, expiresAt)

		assert.Error(t, err)
	})

	t.Run("rejects an empty id_token", func(t *testing.T) {
		t.Parallel()
		idpId, _, accessToken, claims, expiresAt := validUserIdentityArgs()

		_, err := sessionvos.NewUserIdentity(idpId, "", accessToken, nil, claims, expiresAt)

		assert.Error(t, err)
	})

	t.Run("rejects an empty access_token", func(t *testing.T) {
		t.Parallel()
		idpId, idToken, _, claims, expiresAt := validUserIdentityArgs()

		_, err := sessionvos.NewUserIdentity(idpId, idToken, "", nil, claims, expiresAt)

		assert.Error(t, err)
	})

	t.Run("rejects an empty claims.sub", func(t *testing.T) {
		t.Parallel()
		idpId, idToken, accessToken, claims, expiresAt := validUserIdentityArgs()
		claims.Sub = ""

		_, err := sessionvos.NewUserIdentity(idpId, idToken, accessToken, nil, claims, expiresAt)

		assert.Error(t, err)
	})

	t.Run("rejects an empty claims.email", func(t *testing.T) {
		t.Parallel()
		idpId, idToken, accessToken, claims, expiresAt := validUserIdentityArgs()
		claims.Email = ""

		_, err := sessionvos.NewUserIdentity(idpId, idToken, accessToken, nil, claims, expiresAt)

		assert.Error(t, err)
	})
}
