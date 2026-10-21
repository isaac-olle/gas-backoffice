package admins

import (
	"testing"

	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/admins/vos"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewAdmin(t *testing.T) {
	t.Parallel()

	t.Run("creates admin with valid data and records domain event", func(t *testing.T) {
		t.Parallel()
		adminId := ids2.NewUserId()

		admin, err := NewAdmin(adminId, "My Admin", "admin@example.com", false)
		require.NoError(t, err)

		prims := admin.ToPrimitives()
		assert.Equal(t, adminId.String(), prims.Id)
		assert.Equal(t, "My Admin", prims.Name)
		assert.Equal(t, "admin@example.com", prims.Email)
		assert.Equal(t, 1, prims.Version)

		events := admin.PullDomainEvents()
		require.Len(t, events, 1)
		assert.Equal(t, "backoffice.admin.created", events[0].Type().String())
	})

	t.Run("rejects malformed email", func(t *testing.T) {
		t.Parallel()
		adminId := ids2.NewUserId()
		_, err := NewAdmin(adminId, "My Admin", "not-an-email", false)
		assert.Error(t, err)
	})

	t.Run("rejects empty name", func(t *testing.T) {
		t.Parallel()
		adminId := ids2.NewUserId()
		_, err := NewAdmin(adminId, "", "", false)
		assert.Error(t, err)
	})

	t.Run("creates system admin when isSystem is true", func(t *testing.T) {
		t.Parallel()
		adminId := ids2.NewUserId()
		admin, err := NewAdmin(adminId, "System Admin", "", true)
		require.NoError(t, err)
		assert.True(t, admin.IsSystemResource())
	})
}

func TestFromPrimitives(t *testing.T) {
	t.Parallel()

	t.Run("reconstructs admin without recording events", func(t *testing.T) {
		t.Parallel()
		adminId := ids2.NewUserId()

		prims := AdminPrimitives{
			Id:    adminId.String(),
			Name:  "My Admin",
			Email: "admin@example.com",
		}

		admin, err := FromPrimitives(prims)
		require.NoError(t, err)

		got := admin.ToPrimitives()
		assert.Equal(t, adminId.String(), got.Id)
		assert.Equal(t, "My Admin", got.Name)
		assert.Equal(t, "admin@example.com", got.Email)
		assert.Equal(t, 1, got.Version)
		assert.Empty(t, admin.PullDomainEvents())
	})

	t.Run("preserves a non-default version", func(t *testing.T) {
		t.Parallel()
		prims := AdminPrimitives{
			Id:      ids2.NewUserId().String(),
			Name:    "My Admin",
			Version: 5,
		}

		admin, err := FromPrimitives(prims)
		require.NoError(t, err)
		assert.Equal(t, 5, admin.ToPrimitives().Version)
	})

	t.Run("rejects invalid admin id", func(t *testing.T) {
		t.Parallel()
		_, err := FromPrimitives(AdminPrimitives{Id: "bad-id", Name: "Test"})
		assert.Error(t, err)
	})
}

func TestAdmin_EnsureSelf(t *testing.T) {
	t.Parallel()

	t.Run("returns nil when caller is the admin itself", func(t *testing.T) {
		t.Parallel()
		adminId := ids2.NewUserId()
		admin, err := NewAdmin(adminId, "A", "", false)
		require.NoError(t, err)

		assert.NoError(t, admin.EnsureSelf(adminId))
	})

	t.Run("returns error when caller is not the admin", func(t *testing.T) {
		t.Parallel()
		admin, err := NewAdmin(ids2.NewUserId(), "A", "", false)
		require.NoError(t, err)

		assert.Error(t, admin.EnsureSelf(ids2.NewUserId()))
	})
}

func TestAdmin_ChangeName(t *testing.T) {
	t.Parallel()

	t.Run("updates name and records domain event", func(t *testing.T) {
		t.Parallel()
		adminId := ids2.NewUserId()
		admin, err := NewAdmin(adminId, "Old Name", "", false)
		require.NoError(t, err)
		admin.PullDomainEvents() // clear creation event

		newName, err := vos.NewUserName("New Name")
		require.NoError(t, err)

		require.NoError(t, admin.ChangeName(newName))
		assert.Equal(t, "New Name", admin.ToPrimitives().Name)

		events := admin.PullDomainEvents()
		require.Len(t, events, 1)
		assert.Equal(t, "backoffice.admin.name_changed", events[0].Type().String())
	})
}

func TestAdmin_IsSystemResource(t *testing.T) {
	t.Parallel()

	t.Run("returns false by default", func(t *testing.T) {
		t.Parallel()
		admin, err := NewAdmin(ids2.NewUserId(), "A", "", false)
		require.NoError(t, err)
		assert.False(t, admin.IsSystemResource())
	})

	t.Run("returns true when set via FromPrimitives", func(t *testing.T) {
		t.Parallel()
		prims := AdminPrimitives{
			Id:       ids2.NewUserId().String(),
			Name:     "A",
			IsSystem: true,
		}
		admin, err := FromPrimitives(prims)
		require.NoError(t, err)
		assert.True(t, admin.IsSystemResource())
	})
}
