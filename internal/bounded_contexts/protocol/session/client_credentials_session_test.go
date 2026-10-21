package session_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

func TestNewClientCredentialsSession(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		clientID   string
		userId     string
		scope      string
		wantErr    bool
		wantEvents int
	}{
		{
			name:       "creates session and records domain event",
			clientID:   uuid.New().String(),
			userId:     uuid.New().String(),
			scope:      "openid profile",
			wantEvents: 1,
		},
		{
			name:     "rejects invalid client id",
			clientID: "bad-id",
			userId:   uuid.New().String(),
			scope:    "openid",
			wantErr:  true,
		},
		{
			name:     "rejects invalid user id",
			clientID: uuid.New().String(),
			userId:   "bad-user",
			scope:    "openid",
			wantErr:  true,
		},
		{
			name:     "rejects empty scope",
			clientID: uuid.New().String(),
			userId:   uuid.New().String(),
			scope:    "",
			wantErr:  true,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			sess, err := session.NewClientCredentialsSession(ids.NewSessionId(), tc.clientID, tc.userId, tc.scope, "")
			if tc.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, sess)
			assert.Len(t, sess.PullDomainEvents(), tc.wantEvents)
		})
	}
}

func TestClientCredentialsSessionFromPrimitives(t *testing.T) {
	t.Parallel()

	t.Run("reconstructs session from primitives", func(t *testing.T) {
		t.Parallel()
		p := &session.ClientCredentialsSessionPrimitives{
			BaseSessionPrimitives: &session.BaseSessionPrimitives{
				ClientID: uuid.New().String(),
				Scope:    "read",
			},
			ID:       uuid.New().String(),
			UserId:   uuid.New().String(),
			IssuedAt: time.Now().UTC().Format(time.RFC3339),
		}

		sess, err := session.ClientCredentialsSessionFromPrimitives(p)
		require.NoError(t, err)
		require.NotNil(t, sess)

		prim := sess.ToPrimitives()
		assert.Equal(t, p.ClientID, prim.ClientID)
		assert.Equal(t, p.UserId, prim.UserId)
		assert.Equal(t, p.Scope, prim.Scope)
	})
}

func TestClientCredentialsSessionToPrimitives_RoundTrip(t *testing.T) {
	t.Parallel()

	t.Run("primitives round-trip preserves ids and scope", func(t *testing.T) {
		t.Parallel()
		clientID := uuid.New().String()
		userId := uuid.New().String()

		orig, err := session.NewClientCredentialsSession(ids.NewSessionId(), clientID, userId, "openid", "")
		require.NoError(t, err)

		prim := orig.ToPrimitives()
		assert.Equal(t, clientID, prim.ClientID)
		assert.Equal(t, userId, prim.UserId)
	})
}
