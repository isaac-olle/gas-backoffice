package scopes

import (
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewScope(t *testing.T) {
	t.Parallel()

	t.Run("creates scope with valid data and records domain event", func(t *testing.T) {
		t.Parallel()
		scopeId := ids2.NewScopeId()
		prId := ids2.NewProtectedResourceId()
		userId := ids2.NewUserId()

		scope, err := NewScope(scopeId, "read:users", NewOwnershipFromVos(prId, ids2.NewUserId(), userId), false)
		require.NoError(t, err)

		prims := scope.ToPrimitives()
		assert.Equal(t, scopeId.String(), prims.Id)
		assert.Equal(t, "read:users", prims.Name)
		assert.Equal(t, prId.String(), prims.ProtectedResourceId)
		assert.Equal(t, 1, prims.Version)

		events := scope.PullDomainEvents()
		require.Len(t, events, 1)
		assert.Equal(t, "backoffice.scope.created", events[0].Type().String())
	})

	t.Run("rejects empty name", func(t *testing.T) {
		t.Parallel()
		_, err := NewScope(
			ids2.NewScopeId(),
			"",
			NewOwnershipFromVos(ids2.NewProtectedResourceId(), ids2.NewUserId(), ids2.NewUserId()),
			false,
		)
		assert.Error(t, err)
	})
}

func TestScopeFromPrimitives(t *testing.T) {
	t.Parallel()

	t.Run("reconstructs scope without recording events", func(t *testing.T) {
		t.Parallel()
		scopeId := ids2.NewScopeId()
		prId := ids2.NewProtectedResourceId()
		userId := ids2.NewUserId()
		clientId := ids2.NewClientId()

		prims := ScopePrimitives{
			Id:                  scopeId.String(),
			Name:                "read:users",
			ProtectedResourceId: prId.String(),
			AdminId:             ids2.NewUserId().String(),
			UserId:              userId.String(),
			AllowedClientIds:    []string{clientId.String()},
		}

		scope, err := FromPrimitives(prims)
		require.NoError(t, err)

		got := scope.ToPrimitives()
		assert.Equal(t, scopeId.String(), got.Id)
		assert.Equal(t, "read:users", got.Name)
		assert.Len(t, got.AllowedClientIds, 1)
		assert.Equal(t, 1, got.Version)
		assert.Empty(t, scope.PullDomainEvents())
	})

	t.Run("preserves a non-default version", func(t *testing.T) {
		t.Parallel()
		prims := ScopePrimitives{
			Id:                  ids2.NewScopeId().String(),
			Name:                "read:users",
			ProtectedResourceId: ids2.NewProtectedResourceId().String(),
			AdminId:             ids2.NewUserId().String(),
			UserId:              ids2.NewUserId().String(),
			Version:             5,
		}

		scope, err := FromPrimitives(prims)
		require.NoError(t, err)
		assert.Equal(t, 5, scope.ToPrimitives().Version)
	})

	invalidCases := []struct {
		desc  string
		prims ScopePrimitives
	}{
		{
			"rejects invalid scope id",
			ScopePrimitives{
				Id:                  "bad-id",
				Name:                "read:users",
				ProtectedResourceId: ids2.NewProtectedResourceId().String(),
				AdminId:             ids2.NewUserId().String(),
				UserId:              ids2.NewUserId().String(),
			},
		},
		{
			"rejects invalid protected resource id",
			ScopePrimitives{
				Id:                  ids2.NewScopeId().String(),
				Name:                "read:users",
				ProtectedResourceId: "bad-pr-id",
				AdminId:             ids2.NewUserId().String(),
				UserId:              ids2.NewUserId().String(),
			},
		},
		{
			"rejects invalid admin id",
			ScopePrimitives{
				Id:                  ids2.NewScopeId().String(),
				Name:                "read:users",
				ProtectedResourceId: ids2.NewProtectedResourceId().String(),
				AdminId:             "bad-admin-id",
				UserId:              ids2.NewUserId().String(),
			},
		},
	}
	for _, tc := range invalidCases {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()
			_, err := FromPrimitives(tc.prims)
			assert.Error(t, err)
		})
	}
}

func TestScope_EnsureOwnership(t *testing.T) {
	t.Parallel()

	prId := ids2.NewProtectedResourceId()
	adminId := ids2.NewUserId()
	ownerId := ids2.NewUserId()
	newOwnedScope := func(t *testing.T) *Scope {
		t.Helper()
		scope, err := NewScope(ids2.NewScopeId(), "read:users", NewOwnershipFromVos(prId, adminId, ownerId), false)
		require.NoError(t, err)
		return scope
	}

	t.Run("returns nil when ownership matches", func(t *testing.T) {
		t.Parallel()
		assert.NoError(t, newOwnedScope(t).EnsureOwnership(newOwnership(prId, adminId, ownerId)))
	})

	t.Run("returns nil when the caller is the owning admin", func(t *testing.T) {
		t.Parallel()
		assert.NoError(t, newOwnedScope(t).EnsureOwnership(newOwnership(prId, adminId, adminId)))
	})

	t.Run("returns error when protected resource id does not match", func(t *testing.T) {
		t.Parallel()
		assert.Error(t, newOwnedScope(t).EnsureOwnership(newOwnership(ids2.NewProtectedResourceId(), adminId, ownerId)))
	})

	t.Run("returns error when the owning admin targets another protected resource", func(t *testing.T) {
		t.Parallel()
		assert.Error(t, newOwnedScope(t).EnsureOwnership(newOwnership(ids2.NewProtectedResourceId(), adminId, adminId)))
	})

	t.Run("returns error when the caller is another admin", func(t *testing.T) {
		t.Parallel()
		otherAdmin := ids2.NewUserId()
		assert.Error(t, newOwnedScope(t).EnsureOwnership(newOwnership(prId, otherAdmin, otherAdmin)))
	})

	t.Run("returns error when the caller is a sibling tenant of the same admin", func(t *testing.T) {
		t.Parallel()
		assert.Error(t, newOwnedScope(t).EnsureOwnership(newOwnership(prId, adminId, ids2.NewUserId())))
	})
}

func TestScope_AddAllowedClient(t *testing.T) {
	t.Parallel()

	t.Run("adds client id and records event", func(t *testing.T) {
		t.Parallel()
		scope, err := NewScope(
			ids2.NewScopeId(),
			"read:users",
			NewOwnershipFromVos(ids2.NewProtectedResourceId(), ids2.NewUserId(), ids2.NewUserId()),
			false,
		)
		require.NoError(t, err)
		scope.PullDomainEvents() // clear creation event

		clientId := ids2.NewClientId()
		scope.AddAllowedClient(clientId)

		assert.Len(t, scope.ToPrimitives().AllowedClientIds, 1)

		events := scope.PullDomainEvents()
		require.Len(t, events, 1)
		assert.Equal(t, "backoffice.scope.client_added", events[0].Type().String())
	})
}

func TestScope_IsSystemResource(t *testing.T) {
	t.Parallel()

	t.Run("returns false by default on new scope", func(t *testing.T) {
		t.Parallel()
		scope, err := NewScope(ids2.NewScopeId(), "r:x", NewOwnershipFromVos(ids2.NewProtectedResourceId(), ids2.NewUserId(), ids2.NewUserId()), false)
		require.NoError(t, err)
		assert.False(t, scope.IsSystemResource())
	})

	t.Run("returns true when set via FromPrimitives", func(t *testing.T) {
		t.Parallel()
		prims := ScopePrimitives{
			Id:                  ids2.NewScopeId().String(),
			Name:                "r:x",
			ProtectedResourceId: ids2.NewProtectedResourceId().String(),
			AdminId:             ids2.NewUserId().String(),
			UserId:              ids2.NewUserId().String(),
			IsSystem:            true,
		}
		scope, err := FromPrimitives(prims)
		require.NoError(t, err)
		assert.True(t, scope.IsSystemResource())
	})
}

func TestScope_Delete(t *testing.T) {
	t.Parallel()

	t.Run("records a scope deleted event", func(t *testing.T) {
		t.Parallel()
		scope, err := NewScope(ids2.NewScopeId(), "r:x", NewOwnershipFromVos(ids2.NewProtectedResourceId(), ids2.NewUserId(), ids2.NewUserId()), false)
		require.NoError(t, err)
		scope.PullDomainEvents()

		scope.Delete()

		events := scope.PullDomainEvents()
		require.Len(t, events, 1)
		assert.Equal(t, "backoffice.scope.deleted", events[0].Type().String())
	})
}

func TestOwnershipEquals(t *testing.T) {
	t.Parallel()

	prId := ids2.NewProtectedResourceId()
	adminId := ids2.NewUserId()
	userId := ids2.NewUserId()

	t.Run("returns true for same ownership", func(t *testing.T) {
		t.Parallel()
		o1 := newOwnership(prId, adminId, userId)
		o2 := newOwnership(prId, adminId, userId)
		assert.True(t, o1.Equals(o2))
	})

	t.Run("returns false when protected resource differs", func(t *testing.T) {
		t.Parallel()
		o1 := newOwnership(prId, adminId, userId)
		o2 := newOwnership(ids2.NewProtectedResourceId(), adminId, userId)
		assert.False(t, o1.Equals(o2))
	})
}

func TestNewOwnershipFromStrings(t *testing.T) {
	t.Parallel()

	prId := ids2.NewProtectedResourceId()
	adminId := ids2.NewUserId()
	userId := ids2.NewUserId()

	t.Run("creates ownership from valid string ids", func(t *testing.T) {
		t.Parallel()
		o, err := NewOwnership(prId.String(), adminId.String(), userId.String())
		require.NoError(t, err)
		assert.Equal(t, prId.String(), o.ProtectedResourceId.String())
		assert.Equal(t, adminId.String(), o.AdminId.String())
	})

	t.Run("rejects invalid protected resource id", func(t *testing.T) {
		t.Parallel()
		_, err := NewOwnership("bad-id", adminId.String(), userId.String())
		assert.Error(t, err)
	})

	t.Run("rejects invalid admin id", func(t *testing.T) {
		t.Parallel()
		_, err := NewOwnership(prId.String(), "bad-id", userId.String())
		assert.Error(t, err)
	})

	t.Run("rejects invalid user id", func(t *testing.T) {
		t.Parallel()
		_, err := NewOwnership(prId.String(), adminId.String(), "bad-id")
		assert.Error(t, err)
	})
}

func TestOwnershipToPrimitives(t *testing.T) {
	t.Parallel()
	prId := ids2.NewProtectedResourceId()
	adminId := ids2.NewUserId()
	userId := ids2.NewUserId()
	o := newOwnership(prId, adminId, userId)
	p := o.ToPrimitives()
	assert.Equal(t, prId.String(), p.ProtectedResourceId)
	assert.Equal(t, adminId.String(), p.AdminId)
	assert.Equal(t, userId.String(), p.UserId)
}

func TestScope_RemoveAllowedClient(t *testing.T) {
	t.Parallel()

	t.Run("removes client id and records event", func(t *testing.T) {
		t.Parallel()
		clientId := ids2.NewClientId()
		prims := ScopePrimitives{
			Id:                  ids2.NewScopeId().String(),
			Name:                "read:users",
			ProtectedResourceId: ids2.NewProtectedResourceId().String(),
			AdminId:             ids2.NewUserId().String(),
			UserId:              ids2.NewUserId().String(),
			AllowedClientIds:    []string{clientId.String()},
		}
		scope, err := FromPrimitives(prims)
		require.NoError(t, err)

		scope.RemoveAllowedClient(clientId)
		assert.Empty(t, scope.ToPrimitives().AllowedClientIds)

		events := scope.PullDomainEvents()
		require.Len(t, events, 1)
		assert.Equal(t, "backoffice.scope.client_removed", events[0].Type().String())
	})
}
