package users

import (
	"testing"

	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/users/vos"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/authors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewUser(t *testing.T) {
	t.Parallel()

	t.Run("creates user with valid data and records domain event", func(t *testing.T) {
		t.Parallel()
		userId := ids2.NewUserId()

		user, err := NewUser(userId, "My User", "tenant", ids2.NewUserId())
		require.NoError(t, err)

		prims := user.ToPrimitives()
		assert.Equal(t, userId.String(), prims.Id)
		assert.Equal(t, "My User", prims.Name)

		events := user.PullDomainEvents()
		require.Len(t, events, 1)
		assert.Equal(t, "backoffice.user.created", events[0].Type().String())
	})

	invalidNameCases := []struct {
		desc string
		name string
	}{
		{"rejects empty name", ""},
		{"rejects reserved name Gas-User", "Gas-User"},
	}
	for _, tc := range invalidNameCases {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()
			userId := ids2.NewUserId()
			_, err := NewUser(userId, tc.name, "tenant", ids2.NewUserId())
			assert.Error(t, err)
		})
	}
}

func TestFromPrimitives(t *testing.T) {
	t.Parallel()

	t.Run("reconstructs user without recording events", func(t *testing.T) {
		t.Parallel()
		userId := ids2.NewUserId()
		clientId := ids2.NewClientId()
		prId := ids2.NewProtectedResourceId()

		prims := UserPrimitives{
			Id:                   userId.String(),
			AdminId:              ids2.NewUserId().String(),
			Name:                 "My User",
			Role:                 "tenant",
			ClientIds:            []string{clientId.String()},
			ProtectedResourceIds: []string{prId.String()},
		}

		user, err := FromPrimitives(prims)
		require.NoError(t, err)

		got := user.ToPrimitives()
		assert.Equal(t, userId.String(), got.Id)
		assert.Equal(t, "My User", got.Name)
		assert.Len(t, got.ClientIds, 1)
		assert.Len(t, got.ProtectedResourceIds, 1)
		assert.Empty(t, user.PullDomainEvents())
	})

	invalidCases := []struct {
		desc  string
		prims UserPrimitives
	}{
		{
			"rejects invalid user id",
			UserPrimitives{Id: "bad-id", AdminId: ids2.NewUserId().String(), Name: "Test", Role: "tenant"},
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

func TestUser_IsSelf(t *testing.T) {
	t.Parallel()

	t.Run("returns true when author id matches user id", func(t *testing.T) {
		t.Parallel()
		userId := ids2.NewUserId()
		user, err := NewUser(userId, "T", "tenant", ids2.NewUserId())
		require.NoError(t, err)

		author := &authors.UserAuthor{Id: userId, Role: vo.TenantRole}
		assert.True(t, user.IsSelf(author))
	})

	t.Run("returns false when author id does not match", func(t *testing.T) {
		t.Parallel()
		userId := ids2.NewUserId()
		user, err := NewUser(userId, "T", "tenant", ids2.NewUserId())
		require.NoError(t, err)

		author := &authors.UserAuthor{Id: ids2.NewUserId(), Role: vo.TenantRole}
		assert.False(t, user.IsSelf(author))
	})
}

func TestUser_IsOwnedByAdmin(t *testing.T) {
	t.Parallel()

	t.Run("returns true when admin id matches", func(t *testing.T) {
		t.Parallel()
		adminId := ids2.NewUserId()
		user, err := NewUser(ids2.NewUserId(), "T", "tenant", adminId)
		require.NoError(t, err)

		admin := &authors.UserAuthor{Id: adminId, Role: vo.AdminRole}
		assert.True(t, user.IsOwnedByAdmin(admin))
	})

	t.Run("returns false when admin id does not match", func(t *testing.T) {
		t.Parallel()
		user, err := NewUser(ids2.NewUserId(), "T", "tenant", ids2.NewUserId())
		require.NoError(t, err)

		admin := &authors.UserAuthor{Id: ids2.NewUserId(), Role: vo.AdminRole}
		assert.False(t, user.IsOwnedByAdmin(admin))
	})
}

func TestUser_ChangeName(t *testing.T) {
	t.Parallel()

	t.Run("updates name and records domain event", func(t *testing.T) {
		t.Parallel()
		userId := ids2.NewUserId()
		user, err := NewUser(userId, "Old Name", "tenant", ids2.NewUserId())
		require.NoError(t, err)
		user.PullDomainEvents() // clear creation event

		newName, err := vos.NewUserName("New Name")
		require.NoError(t, err)

		require.NoError(t, user.ChangeName(newName))
		assert.Equal(t, "New Name", user.ToPrimitives().Name)

		events := user.PullDomainEvents()
		require.Len(t, events, 1)
		assert.Equal(t, "backoffice.user.name_changed", events[0].Type().String())
	})
}

func TestUser_AddClientId(t *testing.T) {
	t.Parallel()

	t.Run("adds client id to user", func(t *testing.T) {
		t.Parallel()
		userId := ids2.NewUserId()
		user, err := NewUser(userId, "T", "tenant", ids2.NewUserId())
		require.NoError(t, err)

		user.AddClientId(ids2.NewClientId())
		assert.Len(t, user.ToPrimitives().ClientIds, 1)
	})
}

func TestUser_IsSystemResource(t *testing.T) {
	t.Parallel()

	t.Run("returns false by default", func(t *testing.T) {
		t.Parallel()
		userId := ids2.NewUserId()
		user, err := NewUser(userId, "T", "tenant", ids2.NewUserId())
		require.NoError(t, err)
		assert.False(t, user.IsSystemResource())
	})

	t.Run("returns true when set via FromPrimitives", func(t *testing.T) {
		t.Parallel()
		prims := UserPrimitives{
			Id:       ids2.NewUserId().String(),
			AdminId:  ids2.NewUserId().String(),
			Name:     "T",
			Role:     "tenant",
			IsSystem: true,
		}
		user, err := FromPrimitives(prims)
		require.NoError(t, err)
		assert.True(t, user.IsSystemResource())
	})
}

func TestUser_RemoveClientId(t *testing.T) {
	t.Parallel()

	t.Run("removes an existing client id", func(t *testing.T) {
		t.Parallel()
		userId := ids2.NewUserId()
		user, err := NewUser(userId, "T", "tenant", ids2.NewUserId())
		require.NoError(t, err)

		clientId := ids2.NewClientId()
		user.AddClientId(clientId)
		user.RemoveClientId(clientId)

		assert.Empty(t, user.ToPrimitives().ClientIds)
	})
}

func TestUser_RemoveProtectedResourceId(t *testing.T) {
	t.Parallel()

	t.Run("removes an existing protected resource id", func(t *testing.T) {
		t.Parallel()
		userId := ids2.NewUserId()
		user, err := NewUser(userId, "T", "tenant", ids2.NewUserId())
		require.NoError(t, err)

		prId := ids2.NewProtectedResourceId()
		user.AddProtectedResourceId(prId)
		user.RemoveProtectedResourceId(prId)

		assert.Empty(t, user.ToPrimitives().ProtectedResourceIds)
	})
}

func TestUser_Delete(t *testing.T) {
	t.Parallel()

	t.Run("records a deleted domain event", func(t *testing.T) {
		t.Parallel()
		userId := ids2.NewUserId()
		user, err := NewUser(userId, "T", "tenant", ids2.NewUserId())
		require.NoError(t, err)
		user.PullDomainEvents()

		user.Delete()

		events := user.PullDomainEvents()
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

func TestUser_AddProtectedResourceId(t *testing.T) {
	t.Parallel()

	t.Run("adds protected resource id to user", func(t *testing.T) {
		t.Parallel()
		userId := ids2.NewUserId()
		user, err := NewUser(userId, "T", "tenant", ids2.NewUserId())
		require.NoError(t, err)

		user.AddProtectedResourceId(ids2.NewProtectedResourceId())
		assert.Len(t, user.ToPrimitives().ProtectedResourceIds, 1)
	})
}
