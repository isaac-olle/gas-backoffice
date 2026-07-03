package infrastructure

import (
	"context"
	event2 "github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"
	shared_logger "github.com/isaac-olle/gas-backoffice/internal/shared/logger"
	"github.com/isaac-olle/gas-backoffice/internal/shared/workers"
	"sync"
)

type InMemoryEventBus struct {
	mu          sync.RWMutex
	subscribers map[event2.Type][]event2.Subscriber
	logger      shared_logger.Logger
	pool        *workers.Pool
}

func NewInMemoryEventBus(logger shared_logger.Logger, pool *workers.Pool, subscribers map[event2.Type][]event2.Subscriber) *InMemoryEventBus {
	return &InMemoryEventBus{
		subscribers: subscribers,
		logger:      logger,
		pool:        pool,
	}
}

func (b *InMemoryEventBus) Register(subscriber event2.Subscriber) {
	b.mu.Lock()
	defer b.mu.Unlock()
	t := subscriber.SubscribedTo()
	b.subscribers[t] = append(b.subscribers[t], subscriber)
}

func (b *InMemoryEventBus) Publish(ctx context.Context, events []event2.Event) {
	for _, evt := range events {
		b.mu.RLock()
		subs := append(b.subscribers[event2.AllTypes], b.subscribers[evt.Type()]...)
		b.mu.RUnlock()

		for _, sub := range subs {
			b.pool.Process(context.WithoutCancel(ctx), func(ctx context.Context) {
				if err := sub.On(ctx, evt); err != nil {
					b.logger.Error("event handler error",
						shared_logger.Field{Key: "event_type", Value: evt.Type().String()},
						shared_logger.Field{Key: "error", Value: err.Error()},
					)
				}
			})
		}
	}
}
