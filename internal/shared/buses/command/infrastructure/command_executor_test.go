package infrastructure

import (
	"context"
	"errors"
	"testing"

	command2 "github.com/isaac-olle/gas-backoffice/internal/shared/buses/command"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type stubUoW struct{ called bool }

func (u *stubUoW) Do(ctx context.Context, fn func(context.Context) error) error {
	u.called = true
	return fn(ctx)
}

type recordingBus struct {
	dispatched []command2.Type
	failOn     command2.Type
	err        error
	registered bool
}

func (b *recordingBus) Dispatch(_ context.Context, cmd command2.Command) error {
	b.dispatched = append(b.dispatched, cmd.Type())
	if b.failOn != "" && cmd.Type() == b.failOn {
		return b.err
	}
	return nil
}
func (b *recordingBus) Register(command2.Type, command2.Handler) { b.registered = true }

type execCmd struct{ t command2.Type }

func (c execCmd) Type() command2.Type { return c.t }

func TestCommandExecutor_SingleExecution_DelegatesToBus(t *testing.T) {
	t.Parallel()
	uow := &stubUoW{}
	bus := &recordingBus{}
	exec := NewCommandExecutor(bus, uow)

	require.NoError(t, exec.SingleExecution(context.Background(), execCmd{"a"}))

	assert.Equal(t, []command2.Type{"a"}, bus.dispatched)
	assert.False(t, uow.called, "single delegates to the bus (whose decorator opens the tx); the executor does not open its own unit of work")
}

func TestCommandExecutor_MultipleExecution(t *testing.T) {
	t.Parallel()
	uow := &stubUoW{}
	bus := &recordingBus{}
	exec := NewCommandExecutor(bus, uow)

	require.NoError(t, exec.MultipleExecution(context.Background(), execCmd{"a"}, execCmd{"b"}, execCmd{"c"}))

	assert.True(t, uow.called)
	assert.Equal(t, []command2.Type{"a", "b", "c"}, bus.dispatched, "all commands dispatched in order within one unit of work")
}

func TestCommandExecutor_MultipleExecution_StopsOnError(t *testing.T) {
	t.Parallel()
	uow := &stubUoW{}
	bus := &recordingBus{failOn: "b", err: errors.New("boom")}
	exec := NewCommandExecutor(bus, uow)

	err := exec.MultipleExecution(context.Background(), execCmd{"a"}, execCmd{"b"}, execCmd{"c"})

	assert.Error(t, err)
	assert.Equal(t, []command2.Type{"a", "b"}, bus.dispatched, "must stop at the failing command")
}
