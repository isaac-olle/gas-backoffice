package tenants

import (
	ids2 "github.com/isaac-olle/gas/internal/shared/domain/vo/ids"
	"testing"

	"github.com/isaac-olle/gas/internal/bounded_contexts/backoffice/tenants/vos"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewTenant(t *testing.T) {
	t.Parallel()

	t.Run("creates tenant with valid data and records domain event", func(t *testing.T) {
		t.Parallel()
		tenantId := ids2.NewTenantId()
		adminId := ids2.NewAdminId()

		tenant, err := NewTenant(tenantId, "My Tenant", NewOwnership(adminId))
		require.NoError(t, err)

		prims := tenant.ToPrimitives()
		assert.Equal(t, tenantId.String(), prims.Id)
		assert.Equal(t, "My Tenant", prims.Name)
		assert.Equal(t, adminId.String(), prims.AdminId)

		events := tenant.PullDomainEvents()
		require.Len(t, events, 1)
		assert.Equal(t, "backoffice.tenant.created", events[0].Type().String())
	})

	invalidNameCases := []struct {
		desc string
		name string
	}{
		{"rejects empty name", ""},
		{"rejects reserved name Gas-Tenant", "Gas-Tenant"},
	}
	for _, tc := range invalidNameCases {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()
			_, err := NewTenant(ids2.NewTenantId(), tc.name, NewOwnership(ids2.NewAdminId()))
			assert.Error(t, err)
		})
	}
}

func TestFromPrimitives(t *testing.T) {
	t.Parallel()

	t.Run("reconstructs tenant without recording events", func(t *testing.T) {
		t.Parallel()
		tenantId := ids2.NewTenantId()
		adminId := ids2.NewAdminId()
		clientId := ids2.NewClientId()
		prId := ids2.NewProtectedResourceId()

		prims := TenantPrimitives{
			Id:                   tenantId.String(),
			Name:                 "My Tenant",
			AdminId:              adminId.String(),
			ClientIds:            []string{clientId.String()},
			ProtectedResourceIds: []string{prId.String()},
		}

		tenant, err := FromPrimitives(prims)
		require.NoError(t, err)

		got := tenant.ToPrimitives()
		assert.Equal(t, tenantId.String(), got.Id)
		assert.Equal(t, "My Tenant", got.Name)
		assert.Len(t, got.ClientIds, 1)
		assert.Len(t, got.ProtectedResourceIds, 1)
		assert.Empty(t, tenant.PullDomainEvents())
	})

	invalidCases := []struct {
		desc  string
		prims TenantPrimitives
	}{
		{
			"rejects invalid tenant id",
			TenantPrimitives{Id: "bad-id", Name: "Test", AdminId: ids2.NewAdminId().String()},
		},
		{
			"rejects invalid admin id",
			TenantPrimitives{Id: ids2.NewTenantId().String(), Name: "Test", AdminId: "bad-admin"},
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

func TestTenant_CheckOwnership(t *testing.T) {
	t.Parallel()

	t.Run("returns true when admin matches", func(t *testing.T) {
		t.Parallel()
		adminId := ids2.NewAdminId()
		tenant, err := NewTenant(ids2.NewTenantId(), "T", NewOwnership(adminId))
		require.NoError(t, err)

		assert.True(t, tenant.CheckOwnership(NewOwnership(adminId)))
	})

	t.Run("returns false when admin does not match", func(t *testing.T) {
		t.Parallel()
		tenant, err := NewTenant(ids2.NewTenantId(), "T", NewOwnership(ids2.NewAdminId()))
		require.NoError(t, err)

		assert.False(t, tenant.CheckOwnership(NewOwnership(ids2.NewAdminId())))
	})
}

func TestTenant_ChangeName(t *testing.T) {
	t.Parallel()

	t.Run("updates name and records domain event", func(t *testing.T) {
		t.Parallel()
		tenant, err := NewTenant(ids2.NewTenantId(), "Old Name", NewOwnership(ids2.NewAdminId()))
		require.NoError(t, err)
		tenant.PullDomainEvents() // clear creation event

		newName, err := vos.NewTenantName("New Name")
		require.NoError(t, err)

		require.NoError(t, tenant.ChangeName(newName))
		assert.Equal(t, "New Name", tenant.ToPrimitives().Name)

		events := tenant.PullDomainEvents()
		require.Len(t, events, 1)
		assert.Equal(t, "backoffice.tenant.name_changed", events[0].Type().String())
	})
}

func TestTenant_AddClientId(t *testing.T) {
	t.Parallel()

	t.Run("adds client id to tenant", func(t *testing.T) {
		t.Parallel()
		tenant, err := NewTenant(ids2.NewTenantId(), "T", NewOwnership(ids2.NewAdminId()))
		require.NoError(t, err)

		tenant.AddClientId(ids2.NewClientId())
		assert.Len(t, tenant.ToPrimitives().ClientIds, 1)
	})
}

func TestTenant_IsSystemResource(t *testing.T) {
	t.Parallel()

	t.Run("returns false by default", func(t *testing.T) {
		t.Parallel()
		tenant, err := NewTenant(ids2.NewTenantId(), "T", NewOwnership(ids2.NewAdminId()))
		require.NoError(t, err)
		assert.False(t, tenant.IsSystemResource())
	})

	t.Run("returns true when set via FromPrimitives", func(t *testing.T) {
		t.Parallel()
		prims := TenantPrimitives{
			Id:       ids2.NewTenantId().String(),
			Name:     "T",
			AdminId:  ids2.NewAdminId().String(),
			IsSystem: true,
		}
		tenant, err := FromPrimitives(prims)
		require.NoError(t, err)
		assert.True(t, tenant.IsSystemResource())
	})
}

func TestTenant_RemoveClientId(t *testing.T) {
	t.Parallel()

	t.Run("removes an existing client id", func(t *testing.T) {
		t.Parallel()
		tenant, err := NewTenant(ids2.NewTenantId(), "T", NewOwnership(ids2.NewAdminId()))
		require.NoError(t, err)

		clientId := ids2.NewClientId()
		tenant.AddClientId(clientId)
		tenant.RemoveClientId(clientId)

		assert.Empty(t, tenant.ToPrimitives().ClientIds)
	})
}

func TestTenant_RemoveProtectedResourceId(t *testing.T) {
	t.Parallel()

	t.Run("removes an existing protected resource id", func(t *testing.T) {
		t.Parallel()
		tenant, err := NewTenant(ids2.NewTenantId(), "T", NewOwnership(ids2.NewAdminId()))
		require.NoError(t, err)

		prId := ids2.NewProtectedResourceId()
		tenant.AddProtectedResourceId(prId)
		tenant.RemoveProtectedResourceId(prId)

		assert.Empty(t, tenant.ToPrimitives().ProtectedResourceIds)
	})
}

func TestTenant_Delete(t *testing.T) {
	t.Parallel()

	t.Run("records a deleted domain event", func(t *testing.T) {
		t.Parallel()
		tenant, err := NewTenant(ids2.NewTenantId(), "T", NewOwnership(ids2.NewAdminId()))
		require.NoError(t, err)
		tenant.PullDomainEvents()

		tenant.Delete()

		events := tenant.PullDomainEvents()
		require.Len(t, events, 1)
		assert.Equal(t, "backoffice.tenant.deleted", events[0].Type().String())
	})
}

func TestOwnershipToPrimitives(t *testing.T) {
	t.Parallel()
	adminId := ids2.NewAdminId()
	o := NewOwnership(adminId)
	p := o.ToPrimitives()
	assert.Equal(t, adminId.String(), p.AdminId)
}

func TestTenant_AddProtectedResourceId(t *testing.T) {
	t.Parallel()

	t.Run("adds protected resource id to tenant", func(t *testing.T) {
		t.Parallel()
		tenant, err := NewTenant(ids2.NewTenantId(), "T", NewOwnership(ids2.NewAdminId()))
		require.NoError(t, err)

		tenant.AddProtectedResourceId(ids2.NewProtectedResourceId())
		assert.Len(t, tenant.ToPrimitives().ProtectedResourceIds, 1)
	})
}
