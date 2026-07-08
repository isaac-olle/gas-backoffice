package aggregate

import (
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"
	sharedtypes "github.com/isaac-olle/gas-backoffice/internal/shared/types"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// stubEvent is a minimal event.Event implementation for testing.
type stubEvent struct {
	event.BaseEvent
}

const stubEventType event.Type = "stub.event"

func (stubEvent) Type() event.Type { return stubEventType }

func newStubEvent(aggregateID string) stubEvent {
	return stubEvent{
		BaseEvent: event.NewEvent(
			aggregateID,
			stubEventType,
			false,
			map[string]any{"key": "value"},
			sharedtypes.AggregateAdmin,
		),
	}
}

func TestNewRoot(t *testing.T) {
	t.Parallel()
	root := NewRoot(Admin)
	require.NotNil(t, root)
	assert.Equal(t, Admin, root.Name())
}

func TestRecord_AddsEvent(t *testing.T) {
	t.Parallel()
	root := NewRoot(Admin)
	evt := newStubEvent("agg-1")
	root.Record(evt)

	events := root.PullDomainEvents()
	require.Len(t, events, 1)
	assert.Equal(t, "agg-1", events[0].AggregateID())
}

func TestRecord_MultipleEvents(t *testing.T) {
	t.Parallel()
	root := NewRoot(User)
	for i := 0; i < 3; i++ {
		root.Record(newStubEvent("agg-multi"))
	}
	events := root.PullDomainEvents()
	assert.Len(t, events, 3)
}

func TestPullDomainEvents_ClearsAfterPull(t *testing.T) {
	t.Parallel()
	root := NewRoot(Client)
	root.Record(newStubEvent("agg-clear"))

	first := root.PullDomainEvents()
	require.Len(t, first, 1)

	second := root.PullDomainEvents()
	assert.Empty(t, second)
}

func TestClear(t *testing.T) {
	t.Parallel()
	root := NewRoot(Scope)
	root.Record(newStubEvent("agg-x"))
	root.Clear()

	events := root.PullDomainEvents()
	assert.Empty(t, events)
}

func TestPullDomainEvents_EventPreservesFields(t *testing.T) {
	t.Parallel()
	root := NewRoot(Admin)
	before := time.Now()
	evt := newStubEvent("my-aggregate-id")
	root.Record(evt)

	events := root.PullDomainEvents()
	require.Len(t, events, 1)

	e := events[0]
	assert.Equal(t, "my-aggregate-id", e.AggregateID())
	assert.Equal(t, event.Type("stub.event"), e.Type())
	assert.False(t, e.External())
	assert.Equal(t, sharedtypes.AggregateAdmin, e.AggregateType())
	assert.False(t, e.OccurredOn().Before(before))
}
