package tenants

import (
	"testing"

	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/tenants/vos"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/authors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewTenant(t *testing.T) {
	t.Parallel()

	t.Run("creates tenant with valid data and records domain event", func(t *testing.T) {
		t.Parallel()
		userId := ids2.NewUserId()

		tenant, err := NewTenant(userId, "My User", ids2.NewUserId(), "", false)
		require.NoError(t, err)

		prims := tenant.ToPrimitives()
		assert.Equal(t, userId.String(), prims.Id)
		assert.Equal(t, "My User", prims.Name)
		assert.Equal(t, "", prims.Email)
		assert.Equal(t, 1, prims.Version)

		events := tenant.PullDomainEvents()
		require.Len(t, events, 1)
		assert.Equal(t, "backoffice.user.created", events[0].Type().String())
	})

	t.Run("creates tenant with a valid email", func(t *testing.T) {
		t.Parallel()
		userId := ids2.NewUserId()

		tenant, err := NewTenant(userId, "My User", ids2.NewUserId(), "user@example.com", false)
		require.NoError(t, err)
		assert.Equal(t, "user@example.com", tenant.ToPrimitives().Email)
	})

	t.Run("rejects malformed email", func(t *testing.T) {
		t.Parallel()
		userId := ids2.NewUserId()
		_, err := NewTenant(userId, "My User", ids2.NewUserId(), "not-an-email", false)
		assert.Error(t, err)
	})

	invalidNameCases := []struct {
		desc string
		name string
	}{
		{"rejects empty name", ""},
	}
	for _, tc := range invalidNameCases {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()
			userId := ids2.NewUserId()
			_, err := NewTenant(userId, tc.name, ids2.NewUserId(), "", false)
			assert.Error(t, err)
		})
	}

	t.Run("creates system tenant when isSystem is true", func(t *testing.T) {
		t.Parallel()
		userId := ids2.NewUserId()
		tenant, err := NewTenant(userId, "System Tenant", ids2.NewUserId(), "", true)
		require.NoError(t, err)
		assert.True(t, tenant.IsSystemResource())
	})
}

func TestFromPrimitives(t *testing.T) {
	t.Parallel()

	t.Run("reconstructs tenant without recording events", func(t *testing.T) {
		t.Parallel()
		userId := ids2.NewUserId()

		prims := TenantPrimitives{
			Id:      userId.String(),
			AdminId: ids2.NewUserId().String(),
			Name:    "My User",
			Email:   "user@example.com",
		}

		tenant, err := FromPrimitives(prims)
		require.NoError(t, err)

		got := tenant.ToPrimitives()
		assert.Equal(t, userId.String(), got.Id)
		assert.Equal(t, "My User", got.Name)
		assert.Equal(t, "user@example.com", got.Email)
		assert.Equal(t, 1, got.Version)
		assert.Empty(t, tenant.PullDomainEvents())
	})

	t.Run("preserves a non-default version", func(t *testing.T) {
		t.Parallel()
		prims := TenantPrimitives{
			Id:      ids2.NewUserId().String(),
			AdminId: ids2.NewUserId().String(),
			Name:    "My User",
			Version: 5,
		}

		tenant, err := FromPrimitives(prims)
		require.NoError(t, err)
		assert.Equal(t, 5, tenant.ToPrimitives().Version)
	})

	invalidCases := []struct {
		desc  string
		prims TenantPrimitives
	}{
		{
			"rejects invalid user id",
			TenantPrimitives{Id: "bad-id", AdminId: ids2.NewUserId().String(), Name: "Test"},
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

func TestTenant_IsSelf(t *testing.T) {
	t.Parallel()

	t.Run("returns true when author id matches tenant id", func(t *testing.T) {
		t.Parallel()
		userId := ids2.NewUserId()
		tenant, err := NewTenant(userId, "T", ids2.NewUserId(), "", false)
		require.NoError(t, err)

		author := &authors.UserAuthor{Id: userId, Role: vo.TenantRole}
		assert.True(t, tenant.IsSelf(author))
	})

	t.Run("returns false when author id does not match", func(t *testing.T) {
		t.Parallel()
		userId := ids2.NewUserId()
		tenant, err := NewTenant(userId, "T", ids2.NewUserId(), "", false)
		require.NoError(t, err)

		author := &authors.UserAuthor{Id: ids2.NewUserId(), Role: vo.TenantRole}
		assert.False(t, tenant.IsSelf(author))
	})
}

func TestTenant_IsOwnedByAdmin(t *testing.T) {
	t.Parallel()

	t.Run("returns true when admin id matches", func(t *testing.T) {
		t.Parallel()
		adminId := ids2.NewUserId()
		tenant, err := NewTenant(ids2.NewUserId(), "T", adminId, "", false)
		require.NoError(t, err)

		admin := &authors.UserAuthor{Id: adminId, Role: vo.AdminRole}
		assert.True(t, tenant.IsOwnedByAdmin(admin))
	})

	t.Run("returns false when admin id does not match", func(t *testing.T) {
		t.Parallel()
		tenant, err := NewTenant(ids2.NewUserId(), "T", ids2.NewUserId(), "", false)
		require.NoError(t, err)

		admin := &authors.UserAuthor{Id: ids2.NewUserId(), Role: vo.AdminRole}
		assert.False(t, tenant.IsOwnedByAdmin(admin))
	})
}

func TestTenant_EnsureTenantAccess(t *testing.T) {
	t.Parallel()

	t.Run("returns nil when admin owns the tenant", func(t *testing.T) {
		t.Parallel()
		adminId := ids2.NewUserId()
		tenant, err := NewTenant(ids2.NewUserId(), "T", adminId, "", false)
		require.NoError(t, err)

		admin := &authors.UserAuthor{Id: adminId, Role: vo.AdminRole}
		assert.NoError(t, tenant.EnsureTenantAccess(admin))
	})

	t.Run("returns error when admin does not own the tenant", func(t *testing.T) {
		t.Parallel()
		tenant, err := NewTenant(ids2.NewUserId(), "T", ids2.NewUserId(), "", false)
		require.NoError(t, err)

		admin := &authors.UserAuthor{Id: ids2.NewUserId(), Role: vo.AdminRole}
		assert.Error(t, tenant.EnsureTenantAccess(admin))
	})

	t.Run("returns nil when tenant accesses itself", func(t *testing.T) {
		t.Parallel()
		userId := ids2.NewUserId()
		tenant, err := NewTenant(userId, "T", ids2.NewUserId(), "", false)
		require.NoError(t, err)

		self := &authors.UserAuthor{Id: userId, Role: vo.TenantRole}
		assert.NoError(t, tenant.EnsureTenantAccess(self))
	})

	t.Run("returns error when tenant accesses another tenant", func(t *testing.T) {
		t.Parallel()
		tenant, err := NewTenant(ids2.NewUserId(), "T", ids2.NewUserId(), "", false)
		require.NoError(t, err)

		other := &authors.UserAuthor{Id: ids2.NewUserId(), Role: vo.TenantRole}
		assert.Error(t, tenant.EnsureTenantAccess(other))
	})
}

func TestTenant_ChangeName(t *testing.T) {
	t.Parallel()

	t.Run("updates name and records domain event", func(t *testing.T) {
		t.Parallel()
		userId := ids2.NewUserId()
		tenant, err := NewTenant(userId, "Old Name", ids2.NewUserId(), "", false)
		require.NoError(t, err)
		tenant.PullDomainEvents() // clear creation event

		newName, err := vos.NewUserName("New Name")
		require.NoError(t, err)

		require.NoError(t, tenant.ChangeName(newName))
		assert.Equal(t, "New Name", tenant.ToPrimitives().Name)

		events := tenant.PullDomainEvents()
		require.Len(t, events, 1)
		assert.Equal(t, "backoffice.user.name_changed", events[0].Type().String())
	})
}

func TestTenant_IsSystemResource(t *testing.T) {
	t.Parallel()

	t.Run("returns false by default", func(t *testing.T) {
		t.Parallel()
		userId := ids2.NewUserId()
		tenant, err := NewTenant(userId, "T", ids2.NewUserId(), "", false)
		require.NoError(t, err)
		assert.False(t, tenant.IsSystemResource())
	})

	t.Run("returns true when set via FromPrimitives", func(t *testing.T) {
		t.Parallel()
		prims := TenantPrimitives{
			Id:       ids2.NewUserId().String(),
			AdminId:  ids2.NewUserId().String(),
			Name:     "T",
			IsSystem: true,
		}
		tenant, err := FromPrimitives(prims)
		require.NoError(t, err)
		assert.True(t, tenant.IsSystemResource())
	})
}

func TestTenant_Delete(t *testing.T) {
	t.Parallel()

	t.Run("records a deleted domain event", func(t *testing.T) {
		t.Parallel()
		userId := ids2.NewUserId()
		tenant, err := NewTenant(userId, "T", ids2.NewUserId(), "", false)
		require.NoError(t, err)
		tenant.PullDomainEvents()

		tenant.Delete()

		events := tenant.PullDomainEvents()
		require.Len(t, events, 1)
		assert.Equal(t, "backoffice.user.deleted", events[0].Type().String())
	})
}

func TestOwnershipToPrimitives(t *testing.T) {
	t.Parallel()
	userId := ids2.NewUserId()
	o := NewOwnership(userId)
	p := o.ToPrimitives()
	assert.Equal(t, userId.String(), p.UserId)
}
