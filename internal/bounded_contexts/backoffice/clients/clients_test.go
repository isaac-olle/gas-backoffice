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
		ownership := clients.NewOwnership(userId)

		client, err := clients.NewClient(
			clientId,
			"My Client",
			ownership,
			[]string{"https://example.com/cb"},
			[]string{"authorization_code"},
		)
		require.NoError(t, err)

		prims := client.ToPrimitives()
		assert.Equal(t, clientId.String(), prims.Id)
		assert.Equal(t, "My Client", prims.Name)
		assert.Equal(t, userId.String(), prims.UserId)
		assert.Len(t, prims.RedirectUris, 1)
		assert.Len(t, prims.AllowedFlows, 1)

		events := client.PullDomainEvents()
		require.Len(t, events, 1)
		assert.Equal(t, "backoffice.client.created", events[0].Type().String())
	})

	errCases := []struct {
		desc         string
		name         string
		redirectUris []string
		flows        []string
	}{
		{"rejects empty name", "", nil, nil},
		{"rejects invalid redirect uri", "My Client", []string{"ftp://bad.com"}, nil},
		{"rejects invalid flow", "My Client", nil, []string{"implicit"}},
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
			)
			assert.Error(t, err)
		})
	}
}

func TestFromPrimitives(t *testing.T) {
	t.Parallel()

	t.Run("reconstructs client without recording events", func(t *testing.T) {
		t.Parallel()
		clientId := ids2.NewClientId()
		userId := ids2.NewUserId()
		scopeId := ids2.NewScopeId()
		keyId := ids2.NewClientAssertionKeyId()
		secretId := ids2.NewClientSecretId()

		prims := clients.ClientPrimitives{
			Id:            clientId.String(),
			Name:          "My Client",
			UserId:        userId.String(),
			AllowedFlows:  []string{"authorization_code"},
			AllowedScopes: []string{scopeId.String()},
			RedirectUris:  []string{"https://example.com/cb"},
			KeyIds:        []string{keyId.String()},
			SecretIds:     []string{secretId.String()},
		}

		client, err := clients.FromPrimitives(prims)
		require.NoError(t, err)

		got := client.ToPrimitives()
		assert.Equal(t, clientId.String(), got.Id)
		assert.Len(t, got.AllowedFlows, 1)
		assert.Len(t, got.AllowedScopes, 1)
		assert.Len(t, got.RedirectUris, 1)
		assert.Len(t, got.KeyIds, 1)
		assert.Len(t, got.SecretIds, 1)

		assert.Empty(t, client.PullDomainEvents())
	})

	t.Run("returns error for invalid id", func(t *testing.T) {
		t.Parallel()
		prims := clients.ClientPrimitives{
			Id:     "bad-id",
			Name:   "C",
			UserId: ids2.NewUserId().String(),
		}
		_, err := clients.FromPrimitives(prims)
		assert.Error(t, err)
	})

	t.Run("returns error for invalid name", func(t *testing.T) {
		t.Parallel()
		prims := clients.ClientPrimitives{
			Id:     ids2.NewClientId().String(),
			Name:   "",
			UserId: ids2.NewUserId().String(),
		}
		_, err := clients.FromPrimitives(prims)
		assert.Error(t, err)
	})

	t.Run("returns error for invalid user id", func(t *testing.T) {
		t.Parallel()
		prims := clients.ClientPrimitives{
			Id:     ids2.NewClientId().String(),
			Name:   "C",
			UserId: "bad-id",
		}
		_, err := clients.FromPrimitives(prims)
		assert.Error(t, err)
	})

	t.Run("returns error for invalid flow", func(t *testing.T) {
		t.Parallel()
		prims := clients.ClientPrimitives{
			Id:           ids2.NewClientId().String(),
			Name:         "C",
			UserId:       ids2.NewUserId().String(),
			AllowedFlows: []string{"implicit"},
		}
		_, err := clients.FromPrimitives(prims)
		assert.Error(t, err)
	})

	t.Run("returns error for invalid scope id", func(t *testing.T) {
		t.Parallel()
		prims := clients.ClientPrimitives{
			Id:            ids2.NewClientId().String(),
			Name:          "C",
			UserId:        ids2.NewUserId().String(),
			AllowedScopes: []string{"not-a-uuid"},
		}
		_, err := clients.FromPrimitives(prims)
		assert.Error(t, err)
	})

	t.Run("returns error for invalid redirect uri", func(t *testing.T) {
		t.Parallel()
		prims := clients.ClientPrimitives{
			Id:           ids2.NewClientId().String(),
			Name:         "C",
			UserId:       ids2.NewUserId().String(),
			RedirectUris: []string{"ftp://not-valid"},
		}
		_, err := clients.FromPrimitives(prims)
		assert.Error(t, err)
	})

	t.Run("returns error for invalid key id", func(t *testing.T) {
		t.Parallel()
		prims := clients.ClientPrimitives{
			Id:     ids2.NewClientId().String(),
			Name:   "C",
			UserId: ids2.NewUserId().String(),
			KeyIds: []string{"not-a-uuid"},
		}
		_, err := clients.FromPrimitives(prims)
		assert.Error(t, err)
	})

	t.Run("returns error for invalid secret id", func(t *testing.T) {
		t.Parallel()
		prims := clients.ClientPrimitives{
			Id:        ids2.NewClientId().String(),
			Name:      "C",
			UserId:    ids2.NewUserId().String(),
			SecretIds: []string{"not-a-uuid"},
		}
		_, err := clients.FromPrimitives(prims)
		assert.Error(t, err)
	})
}

func TestCheckOwnership(t *testing.T) {
	t.Parallel()

	t.Run("returns true when ownership matches", func(t *testing.T) {
		t.Parallel()
		userId := ids2.NewUserId()
		ownership := clients.NewOwnership(userId)

		client, err := clients.NewClient(ids2.NewClientId(), "C", ownership, nil, nil)
		require.NoError(t, err)

		assert.True(t, client.CheckOwnership(ownership))
	})

	t.Run("returns false when user id differs", func(t *testing.T) {
		t.Parallel()
		userId := ids2.NewUserId()
		ownership := clients.NewOwnership(userId)

		client, err := clients.NewClient(ids2.NewClientId(), "C", ownership, nil, nil)
		require.NoError(t, err)

		otherOwnership := clients.NewOwnership(ids2.NewUserId())
		assert.False(t, client.CheckOwnership(otherOwnership))
	})
}

func TestDelete(t *testing.T) {
	t.Parallel()

	t.Run("records deleted domain event", func(t *testing.T) {
		t.Parallel()
		client, err := clients.NewClient(ids2.NewClientId(), "C", aOwnership(t), nil, nil)
		require.NoError(t, err)
		client.PullDomainEvents() // clear creation event

		client.Delete()

		events := client.PullDomainEvents()
		require.Len(t, events, 1)
		assert.Equal(t, "backoffice.client.deleted", events[0].Type().String())
	})
}

func TestAddAllowedFlow(t *testing.T) {
	t.Parallel()

	t.Run("records flow_added domain event", func(t *testing.T) {
		t.Parallel()
		client, err := clients.NewClient(ids2.NewClientId(), "C", aOwnership(t), nil, nil)
		require.NoError(t, err)
		client.PullDomainEvents()

		flow, fErr := vos.NewFlow("pkce")
		require.NoError(t, fErr)
		client.AddAllowedFlow(flow)

		events := client.PullDomainEvents()
		require.Len(t, events, 1)
		assert.Equal(t, "backoffice.client.flow_added", events[0].Type().String())
	})
}

func TestAddRedirectUri(t *testing.T) {
	t.Parallel()

	t.Run("records redirect_uri_added domain event", func(t *testing.T) {
		t.Parallel()
		client, err := clients.NewClient(ids2.NewClientId(), "C", aOwnership(t), nil, nil)
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
		client, err := clients.NewClient(ids2.NewClientId(), "C", aOwnership(t), nil, nil)
		require.NoError(t, err)

		keyId := ids2.NewClientAssertionKeyId()
		client.AddKeyId(keyId)

		assert.Contains(t, client.ToPrimitives().KeyIds, keyId.String())
	})
}

func TestAddSecretId(t *testing.T) {
	t.Parallel()

	t.Run("appends secret id to primitives", func(t *testing.T) {
		t.Parallel()
		client, err := clients.NewClient(ids2.NewClientId(), "C", aOwnership(t), nil, nil)
		require.NoError(t, err)

		secretId := ids2.NewClientSecretId()
		client.AddSecretId(secretId)

		assert.Contains(t, client.ToPrimitives().SecretIds, secretId.String())
	})
}

func TestAddRemoveAllowedScope(t *testing.T) {
	t.Parallel()

	t.Run("adds scope then removes it", func(t *testing.T) {
		t.Parallel()
		client, err := clients.NewClient(ids2.NewClientId(), "C", aOwnership(t), nil, nil)
		require.NoError(t, err)

		scopeId := ids2.NewScopeId()
		client.AddAllowedScope(scopeId)
		assert.Len(t, client.ToPrimitives().AllowedScopes, 1)

		client.RemoveAllowedScope(scopeId)
		assert.Empty(t, client.ToPrimitives().AllowedScopes)
	})
}

func TestOwnershipEquals(t *testing.T) {
	t.Parallel()

	userId := ids2.NewUserId()

	t.Run("returns true for same ids", func(t *testing.T) {
		t.Parallel()
		o1 := clients.NewOwnership(userId)
		o2 := clients.NewOwnership(userId)
		assert.True(t, o1.Equals(o2))
	})

	t.Run("returns false when user id differs", func(t *testing.T) {
		t.Parallel()
		o1 := clients.NewOwnership(userId)
		o2 := clients.NewOwnership(ids2.NewUserId())
		assert.False(t, o1.Equals(o2))
	})
}

func TestOwnershipToPrimitives(t *testing.T) {
	t.Parallel()
	userId := ids2.NewUserId()
	o := clients.NewOwnership(userId)
	p := o.ToPrimitives()
	assert.Equal(t, userId.String(), p.UserId)
}

// ── helpers ───────────────────────────────────────────────────────────────────

func aOwnership(t *testing.T) clients.Ownership {
	t.Helper()
	return clients.NewOwnership(ids2.NewUserId())
}
