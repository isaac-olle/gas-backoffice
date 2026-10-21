package infrastructure_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	verifyAssertion "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/client_assertion_keys/application/verify_assertion"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session/infrastructure"
)

func TestBusClientAssertionVerifier_Verify(t *testing.T) {
	t.Parallel()

	t.Run("dispatches to the query bus and returns its error", func(t *testing.T) {
		t.Parallel()
		bus := &mockQueryBus{err: assert.AnError}
		v := infrastructure.NewBusClientAssertionVerifier(bus)

		_, err := v.Verify(context.Background(), "jwt")

		assert.Error(t, err)
	})

	t.Run("returns an error when the result has an unexpected type", func(t *testing.T) {
		t.Parallel()
		bus := &mockQueryBus{result: "not-claims"}
		v := infrastructure.NewBusClientAssertionVerifier(bus)

		_, err := v.Verify(context.Background(), "jwt")

		assert.Error(t, err)
	})

	t.Run("returns the client id from the claims", func(t *testing.T) {
		t.Parallel()
		bus := &mockQueryBus{result: &verifyAssertion.ClientClaims{ClientID: "client-1"}}
		v := infrastructure.NewBusClientAssertionVerifier(bus)

		clientID, err := v.Verify(context.Background(), "jwt")

		require.NoError(t, err)
		assert.Equal(t, "client-1", clientID)
	})
}
