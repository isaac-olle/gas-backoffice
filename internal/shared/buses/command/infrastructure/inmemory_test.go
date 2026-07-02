package infrastructure

import (
	"context"
	"errors"
	"github.com/isaac-olle/gas/internal/shared/buses/command"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testCmdType  command.Type = "test.command"
	otherCmdType command.Type = "other.command"
)

type testCommand struct{}

func (t testCommand) Type() command.Type { return testCmdType }

type testHandler struct {
	called bool
	err    error
}

func (h *testHandler) Handle(_ context.Context, _ command.Command) error {
	h.called = true
	return h.err
}

func TestCommandBus_RegisterAndDispatch(t *testing.T) {
	t.Parallel()
	bus := NewCommandBus()
	handler := &testHandler{}
	bus.Register(testCmdType, handler)

	err := bus.Dispatch(context.Background(), testCommand{})
	require.NoError(t, err)
	assert.True(t, handler.called)
}

func TestCommandBus_Dispatch_HandlerReturnsError(t *testing.T) {
	t.Parallel()
	bus := NewCommandBus()
	handlerErr := errors.New("handler failure")
	handler := &testHandler{err: handlerErr}
	bus.Register(testCmdType, handler)

	err := bus.Dispatch(context.Background(), testCommand{})
	require.Error(t, err)
	assert.Equal(t, handlerErr, err)
}

func TestCommandBus_Dispatch_UnknownType_ReturnsError(t *testing.T) {
	t.Parallel()
	bus := NewCommandBus()

	err := bus.Dispatch(context.Background(), testCommand{})
	require.Error(t, err, "should error when no handler registered")
	assert.Contains(t, err.Error(), "HANDLER_NOT_FOUND")
}

func TestCommandBus_Register_MultipleHandlers(t *testing.T) {
	t.Parallel()
	bus := NewCommandBus()
	h1 := &testHandler{}
	h2 := &testHandler{}
	bus.Register(testCmdType, h1)
	bus.Register(otherCmdType, h2)

	err := bus.Dispatch(context.Background(), testCommand{})
	require.NoError(t, err)
	assert.True(t, h1.called)
	assert.False(t, h2.called)
}
