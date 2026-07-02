package infrastructure

import (
	"context"
	"errors"
	query2 "github.com/isaac-olle/gas/internal/shared/buses/query"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testQueryType  query2.Type = "test.query"
	otherQueryType query2.Type = "other.query"
)

type testQuery struct{}

func (t testQuery) Type() query2.Type { return testQueryType }

type testQueryHandler struct {
	result any
	err    error
}

func (h *testQueryHandler) Handle(_ context.Context, _ query2.Query) (any, error) {
	return h.result, h.err
}

func TestQueryBus_RegisterAndDispatch(t *testing.T) {
	t.Parallel()
	bus := NewQueryBus()
	handler := &testQueryHandler{result: "response-value"}
	bus.Register(testQueryType, handler)

	res, err := bus.Dispatch(context.Background(), testQuery{})
	require.NoError(t, err)
	assert.Equal(t, "response-value", res)
}

func TestQueryBus_Dispatch_HandlerReturnsError(t *testing.T) {
	t.Parallel()
	bus := NewQueryBus()
	handlerErr := errors.New("query failure")
	handler := &testQueryHandler{err: handlerErr}
	bus.Register(testQueryType, handler)

	_, err := bus.Dispatch(context.Background(), testQuery{})
	require.Error(t, err)
	assert.Equal(t, handlerErr, err)
}

func TestQueryBus_Dispatch_UnknownType_ReturnsError(t *testing.T) {
	t.Parallel()
	bus := NewQueryBus()

	_, err := bus.Dispatch(context.Background(), testQuery{})
	require.Error(t, err)
	assert.True(t, errors.Is(err, ErrHandlerNotFound))
}

func TestQueryBus_Register_MultipleHandlers(t *testing.T) {
	t.Parallel()
	bus := NewQueryBus()
	h1 := &testQueryHandler{result: "from-h1"}
	h2 := &testQueryHandler{result: "from-h2"}
	bus.Register(testQueryType, h1)
	bus.Register(otherQueryType, h2)

	res, err := bus.Dispatch(context.Background(), testQuery{})
	require.NoError(t, err)
	assert.Equal(t, "from-h1", res)
}

func TestDispatchTyped(t *testing.T) {
	t.Parallel()
	bus := NewQueryBus()
	bus.Register(testQueryType, &testQueryHandler{result: "typed-result"})

	res, err := query2.DispatchTyped[string](context.Background(), bus, testQuery{})
	require.NoError(t, err)
	assert.Equal(t, "typed-result", res)
}

func TestDispatchTyped_WrongType_ReturnsError(t *testing.T) {
	t.Parallel()
	bus := NewQueryBus()
	bus.Register(testQueryType, &testQueryHandler{result: 42}) // returns int, not string

	_, err := query2.DispatchTyped[string](context.Background(), bus, testQuery{})
	require.Error(t, err, "should fail when response type doesn't match")
	assert.Contains(t, err.Error(), "INVALID_RESPONSE")
}
