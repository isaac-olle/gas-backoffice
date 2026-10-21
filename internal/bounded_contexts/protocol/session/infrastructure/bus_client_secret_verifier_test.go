package infrastructure_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session/infrastructure"
)

func TestBusClientSecretVerifier_Verify(t *testing.T) {
	t.Parallel()

	t.Run("dispatches to the query bus and returns its error", func(t *testing.T) {
		t.Parallel()
		bus := &mockQueryBus{err: assert.AnError}
		v := infrastructure.NewBusClientSecretVerifier(bus)

		err := v.Verify(context.Background(), "client-1", "secret")

		assert.Error(t, err)
	})

	t.Run("returns nil when the query succeeds", func(t *testing.T) {
		t.Parallel()
		bus := &mockQueryBus{}
		v := infrastructure.NewBusClientSecretVerifier(bus)

		err := v.Verify(context.Background(), "client-1", "secret")

		require.NoError(t, err)
	})
}
