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

func withPrName(name string) prPrimOpt {
	return func(p *protected_resources.ProtectedResourcePrimitives) { p.Name = name }
}

func withPrTenantId(tenantId string) prPrimOpt {
	return func(p *protected_resources.ProtectedResourcePrimitives) { p.TenantId = tenantId }
}

func withPrAdminId(adminId string) prPrimOpt {
	return func(p *protected_resources.ProtectedResourcePrimitives) { p.AdminId = adminId }
}

func withPrScopeIds(scopeIds ...string) prPrimOpt {
	return func(p *protected_resources.ProtectedResourcePrimitives) { p.ScopeIds = scopeIds }
}

func aPrPrimitives(t *testing.T, opts ...prPrimOpt) protected_resources.ProtectedResourcePrimitives {
	t.Helper()
	p := protected_resources.ProtectedResourcePrimitives{
		Id:       ids2.NewProtectedResourceId().String(),
		Name:     "My Resource",
		TenantId: ids2.NewTenantId().String(),
		AdminId:  ids2.NewAdminId().String(),
		ScopeIds: []string{},
	}
	for _, o := range opts {
		o(&p)
	}
	return p
}

func anOwnership(t *testing.T) (ids2.TenantId, ids2.AdminId, protected_resources.Ownership) {
	t.Helper()
	tenantId := ids2.NewTenantId()
	adminId := ids2.NewAdminId()
	return tenantId, adminId, protected_resources.NewOwnership(tenantId, adminId)
}

func aProtectedResource(t *testing.T, name string) *protected_resources.ProtectedResource {
	t.Helper()
	_, _, ownership := anOwnership(t)
	pr, err := protected_resources.NewProtectedResource(ids2.NewProtectedResourceId(), name, ownership)
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
		tenantId, adminId, ownership := anOwnership(t)

		pr, err := protected_resources.NewProtectedResource(prId, "My Resource", ownership)
		require.NoError(t, err)

		prims := pr.ToPrimitives()
		assert.Equal(t, prId.String(), prims.Id)
		assert.Equal(t, "My Resource", prims.Name)
		assert.Equal(t, tenantId.String(), prims.TenantId)
		assert.Equal(t, adminId.String(), prims.AdminId)
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
		{"reserved name returns error", "Gas-Backoffice"},
	}
	for _, tc := range errCases {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()
			_, _, ownership := anOwnership(t)
			_, err := protected_resources.NewProtectedResource(ids2.NewProtectedResourceId(), tc.name, ownership)
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
		tenantId := ids2.NewTenantId()
		adminId := ids2.NewAdminId()
		scopeId := ids2.NewScopeId()

		prims := aPrPrimitives(t,
			withPrId(prId.String()),
			withPrTenantId(tenantId.String()),
			withPrAdminId(adminId.String()),
			withPrScopeIds(scopeId.String()),
		)

		pr, err := protected_resources.FromPrimitives(prims)
		require.NoError(t, err)

		got := pr.ToPrimitives()
		assert.Equal(t, prId.String(), got.Id)
		assert.Contains(t, got.ScopeIds, scopeId.String())
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
		{"invalid tenant id returns error", []prPrimOpt{withPrTenantId("not-a-uuid")}},
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
// AddScopeId / RemoveScopeId
// ---------------------------------------------------------------------------

func TestScopeIds(t *testing.T) {
	t.Parallel()

	t.Run("adds a scope id", func(t *testing.T) {
		t.Parallel()
		pr := aProtectedResource(t, "R")
		scopeId := ids2.NewScopeId()
		pr.AddScopeId(scopeId)
		assert.Contains(t, pr.ToPrimitives().ScopeIds, scopeId.String())
	})

	t.Run("duplicate scope id is not added twice", func(t *testing.T) {
		t.Parallel()
		pr := aProtectedResource(t, "R")
		scopeId := ids2.NewScopeId()
		pr.AddScopeId(scopeId)
		pr.AddScopeId(scopeId)
		assert.Len(t, pr.ToPrimitives().ScopeIds, 1)
	})

	t.Run("removes a scope id", func(t *testing.T) {
		t.Parallel()
		pr := aProtectedResource(t, "R")
		scopeId := ids2.NewScopeId()
		pr.AddScopeId(scopeId)
		pr.RemoveScopeId(scopeId)
		assert.NotContains(t, pr.ToPrimitives().ScopeIds, scopeId.String())
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

	tenantId := ids2.NewTenantId()
	adminId := ids2.NewAdminId()

	t.Run("returns true for same ownership", func(t *testing.T) {
		t.Parallel()
		o1 := protected_resources.NewOwnership(tenantId, adminId)
		o2 := protected_resources.NewOwnership(tenantId, adminId)
		assert.True(t, o1.Equals(o2))
	})

	t.Run("returns false when tenant differs", func(t *testing.T) {
		t.Parallel()
		o1 := protected_resources.NewOwnership(tenantId, adminId)
		o2 := protected_resources.NewOwnership(ids2.NewTenantId(), adminId)
		assert.False(t, o1.Equals(o2))
	})
}

func TestOwnershipToPrimitives(t *testing.T) {
	t.Parallel()
	tenantId := ids2.NewTenantId()
	adminId := ids2.NewAdminId()
	o := protected_resources.NewOwnership(tenantId, adminId)
	p := o.ToPrimitives()
	assert.Equal(t, tenantId.String(), p.TenantId)
	assert.Equal(t, adminId.String(), p.AdminId)
}

// ---------------------------------------------------------------------------
// CheckOwnership
// ---------------------------------------------------------------------------

func TestCheckOwnership(t *testing.T) {
	t.Parallel()

	t.Run("returns true when ownership matches", func(t *testing.T) {
		t.Parallel()
		tenantId, adminId, ownership := anOwnership(t)
		pr, err := protected_resources.NewProtectedResource(ids2.NewProtectedResourceId(), "R", ownership)
		require.NoError(t, err)

		sameOwnership := protected_resources.NewOwnership(tenantId, adminId)
		assert.True(t, pr.CheckOwnership(sameOwnership))
	})

	t.Run("returns false when ownership does not match", func(t *testing.T) {
		t.Parallel()
		_, adminId, ownership := anOwnership(t)
		pr, err := protected_resources.NewProtectedResource(ids2.NewProtectedResourceId(), "R", ownership)
		require.NoError(t, err)

		otherOwnership := protected_resources.NewOwnership(ids2.NewTenantId(), adminId)
		assert.False(t, pr.CheckOwnership(otherOwnership))
	})
}
