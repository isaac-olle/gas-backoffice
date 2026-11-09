package infrastructure_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session/infrastructure"
)

func TestBusIDTokenGenerator_GenerateIDToken(t *testing.T) {
	t.Parallel()

	t.Run("returns error when authenticated user not resolved", func(t *testing.T) {
		t.Parallel()
		gen := infrastructure.NewBusIDTokenGenerator(&mockQueryBus{})

		_, err := gen.GenerateIDToken(context.Background(), newResolvedBaseSession(t), session.NewUserIdentityProjection("", ""))

		require.Error(t, err)
	})

	t.Run("dispatches to query bus when authenticated user is resolved", func(t *testing.T) {
		t.Parallel()
		bus := &mockQueryBus{result: "signed-id-token"}
		gen := infrastructure.NewBusIDTokenGenerator(bus)

		token, err := gen.GenerateIDToken(context.Background(), newResolvedBaseSession(t), session.NewUserIdentityProjection(uuid.New().String(), ""))

		require.NoError(t, err)
		assert.Equal(t, "signed-id-token", token)
	})
}
