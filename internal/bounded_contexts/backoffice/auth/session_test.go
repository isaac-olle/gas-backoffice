package auth_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/auth"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

func TestNewSession(t *testing.T) {
	t.Parallel()

	t.Run("creates the session and records SessionCreated", func(t *testing.T) {
		t.Parallel()
		userId := ids.NewUserId().String()
		s, err := auth.NewSession(userId, "admin", userId)
		require.NoError(t, err)
		assert.Equal(t, userId, s.ToPrimitives().AdminId)
		assert.Len(t, s.PullDomainEvents(), 1)
	})

	t.Run("returns error for an invalid user id", func(t *testing.T) {
		t.Parallel()
		_, err := auth.NewSession("not-a-uuid", "admin", ids.NewUserId().String())
		assert.Error(t, err)
	})

	t.Run("returns error for an invalid admin id", func(t *testing.T) {
		t.Parallel()
		_, err := auth.NewSession(ids.NewUserId().String(), "tenant", "not-a-uuid")
		assert.Error(t, err)
	})

	t.Run("returns error for an invalid role", func(t *testing.T) {
		t.Parallel()
		_, err := auth.NewSession(ids.NewUserId().String(), "nope", ids.NewUserId().String())
		assert.Error(t, err)
	})
}
