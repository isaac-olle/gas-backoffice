package protected_resources_test

import (
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"testing"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/protected_resources"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/protected_resources/vos"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ---------------------------------------------------------------------------
// Object Mother helpers
// ---------------------------------------------------------------------------

type prPrimOpt func(*protected_resources.ProtectedResourcePrimitives)

func withPrId(id string) prPrimOpt {
	return func(p *protected_resources.ProtectedResourcePrimitives) { p.Id = id }
}

func withPrUserId(userId string) prPrimOpt {
	return func(p *protected_resources.ProtectedResourcePrimitives) { p.UserId = userId }
}

func aPrPrimitives(t *testing.T, opts ...prPrimOpt) protected_resources.ProtectedResourcePrimitives {
	t.Helper()
	p := protected_resources.ProtectedResourcePrimitives{
		Id:      ids2.NewProtectedResourceId().String(),
		Name:    "My Resource",
		AdminId: ids2.NewUserId().String(),
		UserId:  ids2.NewUserId().String(),
	}
	for _, o := range opts {
		o(&p)
	}
	return p
}

func anOwnership(t *testing.T) (ids2.UserId, protected_resources.Ownership) {
	t.Helper()
	userId := ids2.NewUserId()
	return userId, protected_resources.NewOwnership(ids2.NewUserId(), userId)
}

func aProtectedResource(t *testing.T, name string) *protected_resources.ProtectedResource {
	t.Helper()
	_, ownership := anOwnership(t)
	pr, err := protected_resources.NewProtectedResource(ids2.NewProtectedResourceId(), name, ownership, false)
	require.NoError(t, err)
	return pr
}

func mustMakeName(t *testing.T, raw string) vos.ProtectedResourceName {
	t.Helper()
	n, err := vos.NewProtectedResourceName(raw)
	require.NoError(t, err)
	return n
}

// ---------------------------------------------------------------------------
// NewProtectedResource
// ---------------------------------------------------------------------------

func TestNewProtectedResource(t *testing.T) {
	t.Parallel()

	t.Run("creates protected resource with expected primitives", func(t *testing.T) {
		t.Parallel()
		prId := ids2.NewProtectedResourceId()
		userId, ownership := anOwnership(t)

		pr, err := protected_resources.NewProtectedResource(prId, "My Resource", ownership, false)
		require.NoError(t, err)

		prims := pr.ToPrimitives()
		assert.Equal(t, prId.String(), prims.Id)
		assert.Equal(t, "My Resource", prims.Name)
		assert.Equal(t, userId.String(), prims.UserId)
		assert.Equal(t, 1, prims.Version)
	})

	t.Run("records a created domain event", func(t *testing.T) {
		t.Parallel()
		pr := aProtectedResource(t, "My Resource")
		events := pr.PullDomainEvents()
		require.Len(t, events, 1)
		assert.Equal(t, "backoffice.protected_resource.created", events[0].Type().String())
	})

	errCases := []struct {
		desc string
		name string
	}{
		{"empty name returns error", ""},
	}
	for _, tc := range errCases {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()
			_, ownership := anOwnership(t)
			_, err := protected_resources.NewProtectedResource(ids2.NewProtectedResourceId(), tc.name, ownership, false)
			assert.Error(t, err)
		})
	}
}

// ---------------------------------------------------------------------------
// FromPrimitives
// ---------------------------------------------------------------------------

func TestFromPrimitives(t *testing.T) {
	t.Parallel()

	t.Run("reconstructs protected resource from valid primitives", func(t *testing.T) {
		t.Parallel()
		prId := ids2.NewProtectedResourceId()
		userId := ids2.NewUserId()

		prims := aPrPrimitives(t,
			withPrId(prId.String()),
			withPrUserId(userId.String()),
		)

		pr, err := protected_resources.FromPrimitives(prims)
		require.NoError(t, err)

		got := pr.ToPrimitives()
		assert.Equal(t, prId.String(), got.Id)
		assert.Equal(t, 1, got.Version)
	})

	t.Run("preserves a non-default version", func(t *testing.T) {
		t.Parallel()
		prims := aPrPrimitives(t)
		prims.Version = 5

		pr, err := protected_resources.FromPrimitives(prims)
		require.NoError(t, err)
		assert.Equal(t, 5, pr.ToPrimitives().Version)
	})

	t.Run("does not record domain events on reconstruction", func(t *testing.T) {
		t.Parallel()
		pr, err := protected_resources.FromPrimitives(aPrPrimitives(t))
		require.NoError(t, err)
		assert.Empty(t, pr.PullDomainEvents())
	})

	errCases := []struct {
		desc string
		opts []prPrimOpt
	}{
		{"invalid id returns error", []prPrimOpt{withPrId("bad-id")}},
		{"invalid user id returns error", []prPrimOpt{withPrUserId("not-a-uuid")}},
	}
	for _, tc := range errCases {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()
			_, err := protected_resources.FromPrimitives(aPrPrimitives(t, tc.opts...))
			assert.Error(t, err)
		})
	}
}

// ---------------------------------------------------------------------------
// UpdateName
// ---------------------------------------------------------------------------

func TestUpdateName(t *testing.T) {
	t.Parallel()

	t.Run("updates name and records name_updated event", func(t *testing.T) {
		t.Parallel()
		pr := aProtectedResource(t, "Old Name")
		pr.PullDomainEvents() // drain creation event

		pr.UpdateName(mustMakeName(t, "New Name"))

		assert.Equal(t, "New Name", pr.ToPrimitives().Name)
		events := pr.PullDomainEvents()
		require.Len(t, events, 1)
		assert.Equal(t, "backoffice.protected_resource.name_updated", events[0].Type().String())
	})
}

// ---------------------------------------------------------------------------
// Delete
// ---------------------------------------------------------------------------

func TestDelete(t *testing.T) {
	t.Parallel()

	t.Run("records a deleted domain event", func(t *testing.T) {
		t.Parallel()
		pr := aProtectedResource(t, "Resource")
		pr.PullDomainEvents() // drain creation event

		pr.Delete()

		events := pr.PullDomainEvents()
		require.Len(t, events, 1)
		assert.Equal(t, "backoffice.protected_resource.deleted", events[0].Type().String())
	})
}

// ---------------------------------------------------------------------------
// IsSystemResource
// ---------------------------------------------------------------------------

func TestIsSystemResource(t *testing.T) {
	t.Parallel()

	t.Run("returns false by default", func(t *testing.T) {
		t.Parallel()
		pr := aProtectedResource(t, "R")
		assert.False(t, pr.IsSystemResource())
	})

	t.Run("returns true when IsSystem is set in primitives", func(t *testing.T) {
		t.Parallel()
		p := aPrPrimitives(t)
		p.IsSystem = true
		pr, err := protected_resources.FromPrimitives(p)
		require.NoError(t, err)
		assert.True(t, pr.IsSystemResource())
	})
}

// ---------------------------------------------------------------------------
// Ownership.Equals / ToPrimitives
// ---------------------------------------------------------------------------

func TestOwnershipEquals(t *testing.T) {
	t.Parallel()

	adminId := ids2.NewUserId()
	userId := ids2.NewUserId()

	t.Run("returns true for same ownership", func(t *testing.T) {
		t.Parallel()
		o1 := protected_resources.NewOwnership(adminId, userId)
		o2 := protected_resources.NewOwnership(adminId, userId)
		assert.True(t, o1.Equals(o2))
	})

	t.Run("returns false when user differs", func(t *testing.T) {
		t.Parallel()
		o1 := protected_resources.NewOwnership(adminId, userId)
		o2 := protected_resources.NewOwnership(adminId, ids2.NewUserId())
		assert.False(t, o1.Equals(o2))
	})
}

func TestOwnershipToPrimitives(t *testing.T) {
	t.Parallel()
	adminId := ids2.NewUserId()
	userId := ids2.NewUserId()
	o := protected_resources.NewOwnership(adminId, userId)
	p := o.ToPrimitives()
	assert.Equal(t, adminId.String(), p.AdminId)
	assert.Equal(t, userId.String(), p.UserId)
}

func TestNewOwnershipFromStrings(t *testing.T) {
	t.Parallel()

	t.Run("creates ownership from valid string ids", func(t *testing.T) {
		t.Parallel()
		adminId := ids2.NewUserId()
		userId := ids2.NewUserId()
		o, err := protected_resources.NewOwnershipFromStrings(adminId.String(), userId.String())
		require.NoError(t, err)
		assert.Equal(t, protected_resources.NewOwnership(adminId, userId), o)
	})

	t.Run("rejects invalid admin id", func(t *testing.T) {
		t.Parallel()
		_, err := protected_resources.NewOwnershipFromStrings("bad-id", ids2.NewUserId().String())
		assert.Error(t, err)
	})

	t.Run("rejects invalid user id", func(t *testing.T) {
		t.Parallel()
		_, err := protected_resources.NewOwnershipFromStrings(ids2.NewUserId().String(), "bad-id")
		assert.Error(t, err)
	})
}

// ---------------------------------------------------------------------------
// EnsureOwnership
// ---------------------------------------------------------------------------

func TestEnsureOwnership(t *testing.T) {
	t.Parallel()

	adminId := ids2.NewUserId()
	ownerId := ids2.NewUserId()
	newOwnedPr := func(t *testing.T) *protected_resources.ProtectedResource {
		t.Helper()
		pr, err := protected_resources.NewProtectedResource(ids2.NewProtectedResourceId(), "R", protected_resources.NewOwnership(adminId, ownerId), false)
		require.NoError(t, err)
		return pr
	}

	t.Run("returns nil when ownership matches", func(t *testing.T) {
		t.Parallel()
		assert.NoError(t, newOwnedPr(t).EnsureOwnership(protected_resources.NewOwnership(adminId, ownerId)))
	})

	t.Run("returns nil when the caller is the owning admin", func(t *testing.T) {
		t.Parallel()
		assert.NoError(t, newOwnedPr(t).EnsureOwnership(protected_resources.NewOwnership(adminId, adminId)))
	})

	t.Run("returns error when ownership does not match", func(t *testing.T) {
		t.Parallel()
		otherOwnership := protected_resources.NewOwnership(ids2.NewUserId(), ids2.NewUserId())
		assert.Error(t, newOwnedPr(t).EnsureOwnership(otherOwnership))
	})

	t.Run("returns error when the caller is another admin", func(t *testing.T) {
		t.Parallel()
		otherAdmin := ids2.NewUserId()
		assert.Error(t, newOwnedPr(t).EnsureOwnership(protected_resources.NewOwnership(otherAdmin, otherAdmin)))
	})

	t.Run("returns error when the caller is a sibling tenant of the same admin", func(t *testing.T) {
		t.Parallel()
		assert.Error(t, newOwnedPr(t).EnsureOwnership(protected_resources.NewOwnership(adminId, ids2.NewUserId())))
	})
}

// ---------------------------------------------------------------------------
// ScopeIds (denormalized, event-driven read cache)
// ---------------------------------------------------------------------------

func TestScopeIds(t *testing.T) {
	t.Parallel()

	t.Run("a new protected resource starts with no scope ids", func(t *testing.T) {
		t.Parallel()
		pr := aProtectedResource(t, "My Resource")
		assert.Empty(t, pr.ToPrimitives().ScopeIds)
	})

	t.Run("AddScopeId and RemoveScopeId update the cache without recording events", func(t *testing.T) {
		t.Parallel()
		pr := aProtectedResource(t, "My Resource")
		pr.PullDomainEvents()
		scopeId := ids2.NewScopeId()

		pr.AddScopeId(scopeId)
		assert.Equal(t, []string{scopeId.String()}, pr.ToPrimitives().ScopeIds)

		pr.RemoveScopeId(scopeId)
		assert.Empty(t, pr.ToPrimitives().ScopeIds)
		assert.Empty(t, pr.PullDomainEvents())
	})

	t.Run("FromPrimitives round-trips scope ids", func(t *testing.T) {
		t.Parallel()
		scopeId := ids2.NewScopeId().String()
		p := aPrPrimitives(t)
		p.ScopeIds = []string{scopeId}

		pr, err := protected_resources.FromPrimitives(p)
		require.NoError(t, err)
		assert.Equal(t, []string{scopeId}, pr.ToPrimitives().ScopeIds)
	})

	t.Run("FromPrimitives rejects an invalid scope id", func(t *testing.T) {
		t.Parallel()
		p := aPrPrimitives(t)
		p.ScopeIds = []string{"not-a-uuid"}

		_, err := protected_resources.FromPrimitives(p)
		assert.Error(t, err)
	})
}
