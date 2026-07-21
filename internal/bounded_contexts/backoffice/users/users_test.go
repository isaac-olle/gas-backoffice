package users

import (
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"testing"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/users/vos"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewUser(t *testing.T) {
	t.Parallel()

	t.Run("creates user with valid data and records domain event", func(t *testing.T) {
		t.Parallel()
		userId := ids2.NewUserId()

		user, err := NewUser(userId, "My User", NewOwnership(userId))
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
			_, err := NewUser(userId, tc.name, NewOwnership(userId))
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
			Name:                 "My User",
			Role:                 "user",
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
			UserPrimitives{Id: "bad-id", Name: "Test", Role: "user"},
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

func TestUser_CheckOwnership(t *testing.T) {
	t.Parallel()

	t.Run("returns true when userId matches", func(t *testing.T) {
		t.Parallel()
		userId := ids2.NewUserId()
		user, err := NewUser(userId, "T", NewOwnership(userId))
		require.NoError(t, err)

		assert.True(t, user.CheckOwnership(NewOwnership(userId)))
	})

	t.Run("returns false when userId does not match", func(t *testing.T) {
		t.Parallel()
		userId := ids2.NewUserId()
		user, err := NewUser(userId, "T", NewOwnership(userId))
		require.NoError(t, err)

		assert.False(t, user.CheckOwnership(NewOwnership(ids2.NewUserId())))
	})
}

func TestUser_ChangeName(t *testing.T) {
	t.Parallel()

	t.Run("updates name and records domain event", func(t *testing.T) {
		t.Parallel()
		userId := ids2.NewUserId()
		user, err := NewUser(userId, "Old Name", NewOwnership(userId))
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
		user, err := NewUser(userId, "T", NewOwnership(userId))
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
		user, err := NewUser(userId, "T", NewOwnership(userId))
		require.NoError(t, err)
		assert.False(t, user.IsSystemResource())
	})

	t.Run("returns true when set via FromPrimitives", func(t *testing.T) {
		t.Parallel()
		prims := UserPrimitives{
			Id:       ids2.NewUserId().String(),
			Name:     "T",
			Role:     "user",
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
		user, err := NewUser(userId, "T", NewOwnership(userId))
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
		user, err := NewUser(userId, "T", NewOwnership(userId))
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
		user, err := NewUser(userId, "T", NewOwnership(userId))
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
		user, err := NewUser(userId, "T", NewOwnership(userId))
		require.NoError(t, err)

		user.AddProtectedResourceId(ids2.NewProtectedResourceId())
		assert.Len(t, user.ToPrimitives().ProtectedResourceIds, 1)
	})
}
