package admins_test

import (
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"testing"

	"github.com/google/uuid"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/admins"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ---------------------------------------------------------------------------
// Object Mother helpers
// ---------------------------------------------------------------------------

type primOpt func(*admins.AdminPrimitives)

func withId(id string) primOpt {
	return func(p *admins.AdminPrimitives) { p.Id = id }
}

func withName(name string) primOpt {
	return func(p *admins.AdminPrimitives) { p.Name = name }
}

func withUserIds(ids ...string) primOpt {
	return func(p *admins.AdminPrimitives) { p.UserIds = ids }
}

func aPrimitives(t *testing.T, opts ...primOpt) admins.AdminPrimitives {
	t.Helper()
	p := admins.AdminPrimitives{
		Id:      ids2.NewAdminId().String(),
		Name:    "Alice",
		UserIds: []string{},
	}
	for _, o := range opts {
		o(&p)
	}
	return p
}

func anAdmin(t *testing.T, opts ...primOpt) *admins.Admin {
	t.Helper()
	a, err := admins.FromPrimitives(aPrimitives(t, opts...))
	require.NoError(t, err)
	return a
}

// ---------------------------------------------------------------------------
// FromPrimitives
// ---------------------------------------------------------------------------

func TestFromPrimitives(t *testing.T) {
	t.Parallel()

	t.Run("reconstructs admin from valid primitives", func(t *testing.T) {
		t.Parallel()
		adminId := ids2.NewAdminId()
		userId := ids2.NewUserId()

		prims := aPrimitives(t,
			withId(adminId.String()),
			withUserIds(userId.String()),
		)

		a, err := admins.FromPrimitives(prims)
		require.NoError(t, err)

		got := a.ToPrimitives()
		assert.Equal(t, adminId.String(), got.Id)
		assert.Equal(t, "Alice", got.Name)
		assert.Len(t, got.UserIds, 1)
		assert.Equal(t, userId.String(), got.UserIds[0])
	})

	t.Run("does not record domain events on reconstruction", func(t *testing.T) {
		t.Parallel()
		a := anAdmin(t)
		assert.Empty(t, a.PullDomainEvents())
	})

	t.Run("accepts empty user ids", func(t *testing.T) {
		t.Parallel()
		a, err := admins.FromPrimitives(aPrimitives(t, withUserIds()))
		require.NoError(t, err)
		assert.Empty(t, a.ToPrimitives().UserIds)
	})

	errCases := []struct {
		desc string
		opts []primOpt
	}{
		{"invalid id returns error", []primOpt{withId("bad-id")}},
		{"invalid user id returns error", []primOpt{withUserIds("bad-user-id")}},
	}
	for _, tc := range errCases {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()
			_, err := admins.FromPrimitives(aPrimitives(t, tc.opts...))
			assert.Error(t, err)
		})
	}
}

// ---------------------------------------------------------------------------
// AddUserId
// ---------------------------------------------------------------------------

func TestAddUserId(t *testing.T) {
	t.Parallel()

	t.Run("adds a new user id", func(t *testing.T) {
		t.Parallel()
		a := anAdmin(t)
		userId := ids2.NewUserId()
		a.AddUserId(userId)
		assert.Contains(t, a.ToPrimitives().UserIds, userId.String())
	})

	t.Run("is idempotent when adding same user twice", func(t *testing.T) {
		t.Parallel()
		userId := ids2.NewUserId()
		a := anAdmin(t, withUserIds(userId.String()))
		a.AddUserId(userId)
		assert.Len(t, a.ToPrimitives().UserIds, 1)
	})
}

// ---------------------------------------------------------------------------
// RemoveUserId
// ---------------------------------------------------------------------------

func TestRemoveUserId(t *testing.T) {
	t.Parallel()

	t.Run("removes an existing user id", func(t *testing.T) {
		t.Parallel()
		userId := ids2.NewUserId()
		a := anAdmin(t, withUserIds(userId.String()))
		a.RemoveUserId(userId)
		assert.Empty(t, a.ToPrimitives().UserIds)
	})

	t.Run("is a no-op when user id does not exist", func(t *testing.T) {
		t.Parallel()
		a := anAdmin(t)
		a.RemoveUserId(ids2.NewUserId())
		assert.Empty(t, a.ToPrimitives().UserIds)
	})
}

// ---------------------------------------------------------------------------
// ToPrimitives round-trip
// ---------------------------------------------------------------------------

func TestToPrimitivesRoundTrip(t *testing.T) {
	t.Parallel()

	t.Run("preserves all fields across serialisation round-trip", func(t *testing.T) {
		t.Parallel()
		adminId := ids2.NewAdminId()
		tid1, tid2 := ids2.NewUserId(), ids2.NewUserId()

		prims := aPrimitives(t,
			withId(adminId.String()),
			withName("Charlie"),
			withUserIds(tid1.String(), tid2.String()),
		)
		a, err := admins.FromPrimitives(prims)
		require.NoError(t, err)

		got := a.ToPrimitives()
		assert.Equal(t, adminId.String(), got.Id)
		assert.Equal(t, "Charlie", got.Name)
		assert.Len(t, got.UserIds, 2)
	})
}

// ---------------------------------------------------------------------------
// Ownership
// ---------------------------------------------------------------------------

func TestOwnership(t *testing.T) {
	t.Parallel()

	t.Run("creates valid ownership from a valid admin id string", func(t *testing.T) {
		t.Parallel()
		adminId := ids2.NewAdminId()
		o, err := admins.NewOwnership(adminId.String())
		require.NoError(t, err)
		assert.Equal(t, adminId.String(), o.AdminId.String())
	})

	t.Run("returns error for invalid admin id string", func(t *testing.T) {
		t.Parallel()
		_, err := admins.NewOwnership("not-a-uuid")
		assert.Error(t, err)
	})

	t.Run("ToPrimitives exposes admin id as string", func(t *testing.T) {
		t.Parallel()
		adminId := ids2.NewAdminId()
		o, err := admins.NewOwnership(adminId.String())
		require.NoError(t, err)
		assert.Equal(t, adminId.String(), o.ToPrimitives().AdminId)
	})
}

// ---------------------------------------------------------------------------
// IsSystemResource
// ---------------------------------------------------------------------------

func TestIsSystemResource(t *testing.T) {
	t.Parallel()

	t.Run("returns false by default", func(t *testing.T) {
		t.Parallel()
		a := anAdmin(t)
		assert.False(t, a.IsSystemResource())
	})

	t.Run("returns true when IsSystem is set in primitives", func(t *testing.T) {
		t.Parallel()
		p := aPrimitives(t)
		p.IsSystem = true
		a, err := admins.FromPrimitives(p)
		require.NoError(t, err)
		assert.True(t, a.IsSystemResource())
	})
}

// ---------------------------------------------------------------------------
// Ownership.Equals
// ---------------------------------------------------------------------------

func TestOwnershipEquals(t *testing.T) {
	t.Parallel()

	t.Run("returns true for same admin id", func(t *testing.T) {
		t.Parallel()
		adminId := ids2.NewAdminId()
		o1, err := admins.NewOwnership(adminId.String())
		require.NoError(t, err)
		o2, err := admins.NewOwnership(adminId.String())
		require.NoError(t, err)
		assert.True(t, o1.Equals(o2))
	})

	t.Run("returns false for different admin id", func(t *testing.T) {
		t.Parallel()
		o1, err := admins.NewOwnership(ids2.NewAdminId().String())
		require.NoError(t, err)
		o2, err := admins.NewOwnership(ids2.NewAdminId().String())
		require.NoError(t, err)
		assert.False(t, o1.Equals(o2))
	})
}

// ---------------------------------------------------------------------------
// AdminId helpers (ids package)
// ---------------------------------------------------------------------------

func TestNewAdminIdFromString(t *testing.T) {
	t.Parallel()

	t.Run("parses a valid UUID string", func(t *testing.T) {
		t.Parallel()
		raw := uuid.New().String()
		id, err := ids2.NewAdminIdFromString(raw)
		require.NoError(t, err)
		assert.Equal(t, raw, id.String())
	})

	t.Run("returns error for invalid UUID string", func(t *testing.T) {
		t.Parallel()
		_, err := ids2.NewAdminIdFromString("not-a-uuid")
		assert.Error(t, err)
	})
}
