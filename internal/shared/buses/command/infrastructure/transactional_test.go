package infrastructure

import (
	"context"
	"testing"

	command2 "github.com/isaac-olle/gas-backoffice/internal/shared/buses/command"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// The decorator delegates the transaction to the UnitOfWork: a Dispatch must run
// inside uow.Do, and the command must reach the inner bus. (Begin/reuse/commit
// behaviour is the UnitOfWork's own test.)
func TestTransactionalCommandBus_WrapsDispatchInUnitOfWork(t *testing.T) {
	t.Parallel()
	uow := &stubUoW{}
	inner := &recordingBus{}
	bus := NewTransactionalCommandBus(inner, uow)

	require.NoError(t, bus.Dispatch(context.Background(), execCmd{"a"}))

	assert.True(t, uow.called, "dispatch must run inside the unit of work")
	assert.Equal(t, []command2.Type{"a"}, inner.dispatched)
}

func TestTransactionalCommandBus_RegisterDelegates(t *testing.T) {
	t.Parallel()
	inner := &recordingBus{}
	bus := NewTransactionalCommandBus(inner, &stubUoW{})
	bus.Register("x", nil)
	assert.True(t, inner.registered)
}
