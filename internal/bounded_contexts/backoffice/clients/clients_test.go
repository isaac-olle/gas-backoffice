package clients_test

import (
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients/vos"
)

func TestNewClient(t *testing.T) {
	t.Parallel()

	t.Run("creates client with valid fields and records domain event", func(t *testing.T) {
		t.Parallel()
		clientId := ids2.NewClientId()
		userId := ids2.NewUserId()
		adminId := ids2.NewUserId()
		ownership := clients.NewOwnership(adminId, userId)

		client, err := clients.NewClient(
			clientId,
			"My Client",
			ownership,
			[]string{"https://example.com/cb"},
			[]string{"authorization_code"},
			true,
			false,
		)
		require.NoError(t, err)

		prims := client.ToPrimitives()
		assert.Equal(t, clientId.String(), prims.Id)
		assert.Equal(t, "My Client", prims.Name)
		assert.Equal(t, adminId.String(), prims.AdminId)
		assert.Equal(t, userId.String(), prims.UserId)
		assert.Len(t, prims.RedirectUris, 1)
		assert.Len(t, prims.AllowedGrantTypes, 1)
		assert.True(t, prims.Public)
		assert.Equal(t, 1, prims.Version)

		events := client.PullDomainEvents()
		require.Len(t, events, 1)
		assert.Equal(t, "backoffice.client.created", events[0].Type().String())
	})

	errCases := []struct {
		desc         string
		name         string
		redirectUris []string
		flows        []string
		public       bool
	}{
		{"rejects empty name", "", nil, nil, false},
		{"rejects invalid redirect uri", "My Client", []string{"ftp://bad.com"}, nil, false},
		{"rejects invalid flow", "My Client", nil, []string{"implicit"}, false},
		{"rejects public client with client_credentials", "My Client", nil, []string{"client_credentials"}, true},
	}
	for _, tc := range errCases {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()
			_, err := clients.NewClient(
				ids2.NewClientId(),
				tc.name,
				aOwnership(t),
				tc.redirectUris,
				tc.flows,
				tc.public,
				false,
			)
			assert.Error(t, err)
		})
	}

	t.Run("accepts public client with authorization_code only", func(t *testing.T) {
		t.Parallel()
		_, err := clients.NewClient(
			ids2.NewClientId(),
			"My Client",
			aOwnership(t),
			nil,
			[]string{"authorization_code"},
			true,
			false,
		)
		require.NoError(t, err)
	})
}

func TestFromPrimitives(t *testing.T) {
	t.Parallel()

	t.Run("reconstructs client without recording events", func(t *testing.T) {
		t.Parallel()
		clientId := ids2.NewClientId()
		userId := ids2.NewUserId()
		keyId := ids2.NewClientAssertionKeyId()
		secretId := ids2.NewClientSecretId()
		scopeId := ids2.NewScopeId()

		prims := clients.ClientPrimitives{
			Id:                clientId.String(),
			Name:              "My Client",
			AdminId:           ids2.NewUserId().String(),
			UserId:            userId.String(),
			AllowedGrantTypes: []string{"authorization_code"},
			RedirectUris:      []string{"https://example.com/cb"},
			AssertionKeyIds:   []string{keyId.String()},
			SecretIds:         []string{secretId.String()},
			AllowedScopes:     []string{scopeId.String()},
		}

		client, err := clients.FromPrimitives(prims)
		require.NoError(t, err)

		got := client.ToPrimitives()
		assert.Equal(t, clientId.String(), got.Id)
		assert.Len(t, got.AllowedGrantTypes, 1)
		assert.Len(t, got.RedirectUris, 1)
		assert.Len(t, got.AssertionKeyIds, 1)
		assert.Len(t, got.SecretIds, 1)
		assert.Equal(t, []string{scopeId.String()}, got.AllowedScopes)
		assert.Equal(t, 1, got.Version)

		assert.Empty(t, client.PullDomainEvents())
	})

	t.Run("defaults version to 1 when primitives omit it", func(t *testing.T) {
		t.Parallel()
		prims := clients.ClientPrimitives{
			Id:      ids2.NewClientId().String(),
			Name:    "My Client",
			AdminId: ids2.NewUserId().String(),
			UserId:  ids2.NewUserId().String(),
		}

		client, err := clients.FromPrimitives(prims)
		require.NoError(t, err)
		assert.Equal(t, 1, client.ToPrimitives().Version)
	})

	t.Run("preserves a non-default version", func(t *testing.T) {
		t.Parallel()
		prims := clients.ClientPrimitives{
			Id:      ids2.NewClientId().String(),
			Name:    "My Client",
			AdminId: ids2.NewUserId().String(),
			UserId:  ids2.NewUserId().String(),
			Version: 5,
		}

		client, err := clients.FromPrimitives(prims)
		require.NoError(t, err)
		assert.Equal(t, 5, client.ToPrimitives().Version)
	})

	t.Run("returns error for invalid id", func(t *testing.T) {
		t.Parallel()
		prims := clients.ClientPrimitives{
			Id:      "bad-id",
			Name:    "C",
			AdminId: ids2.NewUserId().String(),
			UserId:  ids2.NewUserId().String(),
		}
		_, err := clients.FromPrimitives(prims)
		assert.Error(t, err)
	})

	t.Run("returns error for invalid name", func(t *testing.T) {
		t.Parallel()
		prims := clients.ClientPrimitives{
			Id:      ids2.NewClientId().String(),
			Name:    "",
			AdminId: ids2.NewUserId().String(),
			UserId:  ids2.NewUserId().String(),
		}
		_, err := clients.FromPrimitives(prims)
		assert.Error(t, err)
	})

	t.Run("returns error for invalid user id", func(t *testing.T) {
		t.Parallel()
		prims := clients.ClientPrimitives{
			Id:      ids2.NewClientId().String(),
			Name:    "C",
			AdminId: ids2.NewUserId().String(),
			UserId:  "bad-id",
		}
		_, err := clients.FromPrimitives(prims)
		assert.Error(t, err)
	})

	t.Run("returns error for invalid admin id", func(t *testing.T) {
		t.Parallel()
		prims := clients.ClientPrimitives{
			Id:      ids2.NewClientId().String(),
			Name:    "C",
			AdminId: "bad-id",
			UserId:  ids2.NewUserId().String(),
		}
		_, err := clients.FromPrimitives(prims)
		assert.Error(t, err)
	})

	t.Run("returns error for invalid flow", func(t *testing.T) {
		t.Parallel()
		prims := clients.ClientPrimitives{
			Id:                ids2.NewClientId().String(),
			Name:              "C",
			AdminId:           ids2.NewUserId().String(),
			UserId:            ids2.NewUserId().String(),
			AllowedGrantTypes: []string{"implicit"},
		}
		_, err := clients.FromPrimitives(prims)
		assert.Error(t, err)
	})

	t.Run("returns error for invalid redirect uri", func(t *testing.T) {
		t.Parallel()
		prims := clients.ClientPrimitives{
			Id:           ids2.NewClientId().String(),
			Name:         "C",
			AdminId:      ids2.NewUserId().String(),
			UserId:       ids2.NewUserId().String(),
			RedirectUris: []string{"ftp://not-valid"},
		}
		_, err := clients.FromPrimitives(prims)
		assert.Error(t, err)
	})

	t.Run("returns error for invalid key id", func(t *testing.T) {
		t.Parallel()
		prims := clients.ClientPrimitives{
			Id:              ids2.NewClientId().String(),
			Name:            "C",
			AdminId:         ids2.NewUserId().String(),
			UserId:          ids2.NewUserId().String(),
			AssertionKeyIds: []string{"not-a-uuid"},
		}
		_, err := clients.FromPrimitives(prims)
		assert.Error(t, err)
	})

	t.Run("returns error for invalid secret id", func(t *testing.T) {
		t.Parallel()
		prims := clients.ClientPrimitives{
			Id:        ids2.NewClientId().String(),
			Name:      "C",
			AdminId:   ids2.NewUserId().String(),
			UserId:    ids2.NewUserId().String(),
			SecretIds: []string{"not-a-uuid"},
		}
		_, err := clients.FromPrimitives(prims)
		assert.Error(t, err)
	})

	t.Run("returns error for invalid allowed scope id", func(t *testing.T) {
		t.Parallel()
		prims := clients.ClientPrimitives{
			Id:            ids2.NewClientId().String(),
			Name:          "C",
			AdminId:       ids2.NewUserId().String(),
			UserId:        ids2.NewUserId().String(),
			AllowedScopes: []string{"not-a-uuid"},
		}
		_, err := clients.FromPrimitives(prims)
		assert.Error(t, err)
	})
}

func TestEnsureOwnership(t *testing.T) {
	t.Parallel()

	adminId := ids2.NewUserId()
	ownerId := ids2.NewUserId()
	newOwnedClient := func(t *testing.T) *clients.Client {
		t.Helper()
		client, err := clients.NewClient(ids2.NewClientId(), "C", clients.NewOwnership(adminId, ownerId), nil, nil, false, false)
		require.NoError(t, err)
		return client
	}

	t.Run("returns nil when ownership matches", func(t *testing.T) {
		t.Parallel()
		assert.NoError(t, newOwnedClient(t).EnsureOwnership(clients.NewOwnership(adminId, ownerId)))
	})

	t.Run("returns nil when the caller is the owning admin", func(t *testing.T) {
		t.Parallel()
		assert.NoError(t, newOwnedClient(t).EnsureOwnership(clients.NewOwnership(adminId, adminId)))
	})

	t.Run("returns error when user id differs", func(t *testing.T) {
		t.Parallel()
		otherOwnership := clients.NewOwnership(ids2.NewUserId(), ids2.NewUserId())
		assert.Error(t, newOwnedClient(t).EnsureOwnership(otherOwnership))
	})

	t.Run("returns error when the caller is another admin", func(t *testing.T) {
		t.Parallel()
		otherAdmin := ids2.NewUserId()
		assert.Error(t, newOwnedClient(t).EnsureOwnership(clients.NewOwnership(otherAdmin, otherAdmin)))
	})

	t.Run("returns error when the caller is a sibling tenant of the same admin", func(t *testing.T) {
		t.Parallel()
		assert.Error(t, newOwnedClient(t).EnsureOwnership(clients.NewOwnership(adminId, ids2.NewUserId())))
	})
}

func TestDelete(t *testing.T) {
	t.Parallel()

	t.Run("records deleted domain event", func(t *testing.T) {
		t.Parallel()
		client, err := clients.NewClient(ids2.NewClientId(), "C", aOwnership(t), nil, nil, false, false)
		require.NoError(t, err)
		client.PullDomainEvents() // clear creation event

		client.Delete()

		events := client.PullDomainEvents()
		require.Len(t, events, 1)
		assert.Equal(t, "backoffice.client.deleted", events[0].Type().String())
	})
}

func TestAddAllowedGrantType(t *testing.T) {
	t.Parallel()

	t.Run("records grant_type_added domain event", func(t *testing.T) {
		t.Parallel()
		client, err := clients.NewClient(ids2.NewClientId(), "C", aOwnership(t), nil, nil, false, false)
		require.NoError(t, err)
		client.PullDomainEvents()

		grantType, gErr := vos.NewGrantType("client_credentials")
		require.NoError(t, gErr)
		client.AddAllowedGrantType(grantType)

		events := client.PullDomainEvents()
		require.Len(t, events, 1)
		assert.Equal(t, "backoffice.client.grant_type_added", events[0].Type().String())
	})
}

func TestAddRedirectUri(t *testing.T) {
	t.Parallel()

	t.Run("records redirect_uri_added domain event", func(t *testing.T) {
		t.Parallel()
		client, err := clients.NewClient(ids2.NewClientId(), "C", aOwnership(t), nil, nil, false, false)
		require.NoError(t, err)
		client.PullDomainEvents()

		uri, uErr := vos.NewRedirectUri("https://example.com/cb")
		require.NoError(t, uErr)
		client.AddRedirectUri(uri)

		events := client.PullDomainEvents()
		require.Len(t, events, 1)
		assert.Equal(t, "backoffice.client.redirect_uri_added", events[0].Type().String())
	})
}

func TestAddKeyId(t *testing.T) {
	t.Parallel()

	t.Run("appends key id to primitives", func(t *testing.T) {
		t.Parallel()
		client, err := clients.NewClient(ids2.NewClientId(), "C", aOwnership(t), nil, nil, false, false)
		require.NoError(t, err)

		keyId := ids2.NewClientAssertionKeyId()
		client.AddKeyId(keyId)

		assert.Contains(t, client.ToPrimitives().AssertionKeyIds, keyId.String())
	})
}

func TestAddSecretId(t *testing.T) {
	t.Parallel()

	t.Run("appends secret id to primitives", func(t *testing.T) {
		t.Parallel()
		client, err := clients.NewClient(ids2.NewClientId(), "C", aOwnership(t), nil, nil, false, false)
		require.NoError(t, err)

		secretId := ids2.NewClientSecretId()
		client.AddSecretId(secretId)

		assert.Contains(t, client.ToPrimitives().SecretIds, secretId.String())
	})
}

func TestAllowedScopes(t *testing.T) {
	t.Parallel()

	t.Run("new client starts with no allowed scopes", func(t *testing.T) {
		t.Parallel()
		client, err := clients.NewClient(ids2.NewClientId(), "C", aOwnership(t), nil, nil, false, false)
		require.NoError(t, err)

		assert.Empty(t, client.ToPrimitives().AllowedScopes)
	})

	t.Run("adds and removes scope ids without recording events", func(t *testing.T) {
		t.Parallel()
		client, err := clients.NewClient(ids2.NewClientId(), "C", aOwnership(t), nil, nil, false, false)
		require.NoError(t, err)
		client.PullDomainEvents()

		scopeId := ids2.NewScopeId()
		client.AddAllowedScope(scopeId)
		assert.Contains(t, client.ToPrimitives().AllowedScopes, scopeId.String())

		client.RemoveAllowedScope(scopeId)
		assert.NotContains(t, client.ToPrimitives().AllowedScopes, scopeId.String())
		assert.Empty(t, client.PullDomainEvents())
	})
}

func TestOwnershipEquals(t *testing.T) {
	t.Parallel()

	userId := ids2.NewUserId()

	t.Run("returns true for same ids", func(t *testing.T) {
		t.Parallel()
		o1 := clients.NewOwnership(userId, userId)
		o2 := clients.NewOwnership(userId, userId)
		assert.True(t, o1.Equals(o2))
	})

	t.Run("returns false when user id differs", func(t *testing.T) {
		t.Parallel()
		o1 := clients.NewOwnership(userId, userId)
		o2 := clients.NewOwnership(userId, ids2.NewUserId())
		assert.False(t, o1.Equals(o2))
	})
}

func TestOwnershipToPrimitives(t *testing.T) {
	t.Parallel()
	userId := ids2.NewUserId()
	adminId := ids2.NewUserId()
	o := clients.NewOwnership(adminId, userId)
	p := o.ToPrimitives()
	assert.Equal(t, adminId.String(), p.AdminId)
	assert.Equal(t, userId.String(), p.UserId)
}

// ── helpers ───────────────────────────────────────────────────────────────────

func aOwnership(t *testing.T) clients.Ownership {
	t.Helper()
	return clients.NewOwnership(ids2.NewUserId(), ids2.NewUserId())
}
