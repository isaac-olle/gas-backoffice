package auth

import (
	"testing"

	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewSession(t *testing.T) {
	t.Parallel()

	t.Run("creates session with valid data and records domain event", func(t *testing.T) {
		t.Parallel()
		sessionId := ids2.NewSessionId()
		userId := ids2.NewUserId()

		session, err := NewSession(sessionId, userId)
		require.NoError(t, err)

		prims := session.ToPrimitives()
		assert.Equal(t, sessionId.String(), prims.Id)
		assert.Equal(t, userId.String(), prims.UserId)

		events := session.PullDomainEvents()
		require.Len(t, events, 1)
		assert.Equal(t, "backoffice.session.created", events[0].Type().String())
	})
}

func TestFromPrimitives(t *testing.T) {
	t.Parallel()

	t.Run("reconstructs session without recording events", func(t *testing.T) {
		t.Parallel()
		sessionId := ids2.NewSessionId()
		userId := ids2.NewUserId()

		prims := SessionPrimitives{
			Id:     sessionId.String(),
			UserId: userId.String(),
		}

		session, err := FromPrimitives(prims)
		require.NoError(t, err)

		got := session.ToPrimitives()
		assert.Equal(t, sessionId.String(), got.Id)
		assert.Equal(t, userId.String(), got.UserId)
		assert.Empty(t, session.PullDomainEvents())
	})

	invalidCases := []struct {
		desc  string
		prims SessionPrimitives
	}{
		{
			"rejects invalid session id",
			SessionPrimitives{Id: "bad-id", UserId: ids2.NewUserId().String()},
		},
		{
			"rejects invalid user id",
			SessionPrimitives{Id: ids2.NewSessionId().String(), UserId: "bad-id"},
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

func TestSession_BelongsTo(t *testing.T) {
	t.Parallel()

	t.Run("returns true when user id matches", func(t *testing.T) {
		t.Parallel()
		userId := ids2.NewUserId()
		session, err := NewSession(ids2.NewSessionId(), userId)
		require.NoError(t, err)

		assert.True(t, session.BelongsTo(userId))
	})

	t.Run("returns false when user id does not match", func(t *testing.T) {
		t.Parallel()
		session, err := NewSession(ids2.NewSessionId(), ids2.NewUserId())
		require.NoError(t, err)

		assert.False(t, session.BelongsTo(ids2.NewUserId()))
	})
}
