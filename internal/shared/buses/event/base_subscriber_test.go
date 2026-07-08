package event

import (
	"context"
	sharedtypes "github.com/isaac-olle/gas-backoffice/internal/shared/types"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// concreteEvent wraps BaseEvent and overrides Type() so that
// NewBaseSubscriber[concreteEvent] picks up the correct event type
// from the zero value.
type concreteEvent struct {
	BaseEvent
}

const concreteEventType Type = "concrete.event"

func (concreteEvent) Type() Type { return concreteEventType }

func newConcreteEvent(aggID string) concreteEvent {
	return concreteEvent{NewEvent(aggID, concreteEventType, false, nil, sharedtypes.AggregateAdmin)}
}

func TestNewBaseSubscriber_SubscribedTo(t *testing.T) {
	t.Parallel()
	called := false
	sub := NewBaseSubscriber[concreteEvent](func(_ context.Context, _ concreteEvent) error {
		called = true
		return nil
	})

	assert.Equal(t, concreteEventType, sub.SubscribedTo())
	assert.False(t, called)
}

func TestBaseSubscriber_On_CallsHandler(t *testing.T) {
	t.Parallel()
	var received concreteEvent
	sub := NewBaseSubscriber[concreteEvent](func(_ context.Context, e concreteEvent) error {
		received = e
		return nil
	})

	evt := newConcreteEvent("agg-sub")
	err := sub.On(context.Background(), evt)
	require.NoError(t, err)
	assert.Equal(t, "agg-sub", received.AggregateID())
}

func TestBaseSubscriber_On_WrongType_ReturnsNil(t *testing.T) {
	t.Parallel()
	sub := NewBaseSubscriber[concreteEvent](func(_ context.Context, _ concreteEvent) error {
		return nil
	})

	// Pass a raw BaseEvent (not concreteEvent) — type assertion should fail silently.
	rawEvt := NewEvent("agg-wrong", Type("other.event"), false, nil, sharedtypes.AggregateUser)
	err := sub.On(context.Background(), rawEvt)
	assert.NoError(t, err)
}

func TestNewInMemorySubscriber(t *testing.T) {
	t.Parallel()
	sub := NewInMemorySubscriber[concreteEvent](func(_ context.Context, _ concreteEvent) error {
		return nil
	})
	assert.Equal(t, Type("concrete.event"), sub.SubscribedTo())
}
