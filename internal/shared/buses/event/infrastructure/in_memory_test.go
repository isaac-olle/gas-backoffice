package infrastructure

import (
	"context"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	event2 "github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"
	"github.com/isaac-olle/gas-backoffice/internal/shared/logger"
	sharedtypes "github.com/isaac-olle/gas-backoffice/internal/shared/types"
	"github.com/isaac-olle/gas-backoffice/internal/shared/workers"
)

type noopLogger struct{}

func (n *noopLogger) Info(_ string, _ ...logger.Field)  {}
func (n *noopLogger) Error(_ string, _ ...logger.Field) {}
func (n *noopLogger) Fatal(_ string, _ ...logger.Field) {}
func (n *noopLogger) Debug(_ string, _ ...logger.Field) {}

type testEvent struct{ event2.BaseEvent }

const testEventType event2.Type = "test.event"

func (testEvent) Type() event2.Type { return testEventType }

func newTestEvent(aggID string) testEvent {
	return testEvent{event2.NewEvent(aggID, testEventType, false, nil, sharedtypes.AggregateAdmin)}
}

type otherEvent struct{ event2.BaseEvent }

const otherEventType event2.Type = "other.event"

func (otherEvent) Type() event2.Type { return otherEventType }

func newOtherEvent(aggID string) otherEvent {
	return otherEvent{event2.NewEvent(aggID, otherEventType, false, nil, sharedtypes.AggregateTenant)}
}

func makebus() *InMemoryEventBus {
	pool := workers.NewPool(4)
	return NewInMemoryEventBus(&noopLogger{}, pool, make(map[event2.Type][]event2.Subscriber))
}

func waitFor(wg *sync.WaitGroup) {
	done := make(chan struct{})
	go func() { wg.Wait(); close(done) }()
	<-done
}

func TestRegister_AddsSubscriber(t *testing.T) {
	t.Parallel()
	bus := makebus()
	var wg sync.WaitGroup
	wg.Add(1)
	var called bool

	sub := event2.NewInMemorySubscriber[testEvent](func(_ context.Context, _ testEvent) error {
		called = true
		wg.Done()
		return nil
	})
	bus.Register(sub)

	bus.Publish(context.Background(), []event2.Event{newTestEvent("agg-1")})
	waitFor(&wg)

	assert.True(t, called)
}

func TestPublish_DispatchesToMatchingSubscriber(t *testing.T) {
	t.Parallel()
	bus := makebus()
	var mu sync.Mutex
	var received []string
	var wg sync.WaitGroup
	wg.Add(2)

	sub := event2.NewInMemorySubscriber[testEvent](func(_ context.Context, e testEvent) error {
		mu.Lock()
		received = append(received, e.AggregateID())
		mu.Unlock()
		wg.Done()
		return nil
	})
	bus.Register(sub)

	bus.Publish(context.Background(), []event2.Event{
		newTestEvent("agg-a"),
		newTestEvent("agg-b"),
	})
	waitFor(&wg)

	mu.Lock()
	defer mu.Unlock()
	assert.Len(t, received, 2)
	assert.Contains(t, received, "agg-a")
	assert.Contains(t, received, "agg-b")
}

func TestPublish_SkipsNonMatchingSubscribers(t *testing.T) {
	t.Parallel()
	bus := makebus()
	var called bool

	sub := event2.NewInMemorySubscriber[testEvent](func(_ context.Context, _ testEvent) error {
		called = true
		return nil
	})
	bus.Register(sub)

	bus.Publish(context.Background(), []event2.Event{newOtherEvent("agg-other")})

	assert.False(t, called, "subscriber should not be called for a different event type")
}

type allTypesSubscriber struct {
	handler func(context.Context, event2.Event) error
}

func (a *allTypesSubscriber) On(ctx context.Context, e event2.Event) error {
	return a.handler(ctx, e)
}
func (a *allTypesSubscriber) SubscribedTo() event2.Type { return event2.AllTypes }

func TestPublish_AllTypesSubscriberReceivesAll(t *testing.T) {
	t.Parallel()
	bus := makebus()
	var mu sync.Mutex
	var count int
	var wg sync.WaitGroup
	wg.Add(2)

	bus.Register(&allTypesSubscriber{
		handler: func(_ context.Context, _ event2.Event) error {
			mu.Lock()
			count++
			mu.Unlock()
			wg.Done()
			return nil
		},
	})

	bus.Publish(context.Background(), []event2.Event{
		newTestEvent("agg-1"),
		newOtherEvent("agg-2"),
	})
	waitFor(&wg)

	mu.Lock()
	defer mu.Unlock()
	assert.Equal(t, 2, count)
}

func TestPublish_NilEvent_IsHandled(t *testing.T) {
	t.Parallel()
	bus := makebus()
	var called bool
	var wg sync.WaitGroup
	wg.Add(1)

	sub := event2.NewInMemorySubscriber[testEvent](func(_ context.Context, _ testEvent) error {
		called = true
		wg.Done()
		return nil
	})
	bus.Register(sub)

	require.NotPanics(t, func() {
		bus.Publish(context.Background(), []event2.Event{newTestEvent("agg-ok")})
		waitFor(&wg)
	})
	assert.True(t, called)
}
