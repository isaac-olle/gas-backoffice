package infrastructure

import (
	"context"
	event2 "github.com/isaac-olle/gas/internal/shared/buses/event"
	shared_logger "github.com/isaac-olle/gas/internal/shared/logger"
	"time"
)

const defaultBatchSize = 100
const maxConsecutiveErrors = 3

// Consumer polls the outbox table periodically and dispatches pending
// events to subscribers via the InMemoryEventBus.
type Consumer struct {
	outboxRepo        event2.OutboxRepository
	inner             event2.Bus
	registry          *event2.BuilderRegistry
	interval          time.Duration
	batchSize         int
	logger            shared_logger.Logger
	consecutiveErrors int
}

func NewOutboxConsumer(
	outboxRepo event2.OutboxRepository,
	inner event2.Bus,
	registry *event2.BuilderRegistry,
	interval time.Duration,
	logger shared_logger.Logger,
) *Consumer {
	return &Consumer{
		outboxRepo: outboxRepo,
		inner:      inner,
		registry:   registry,
		interval:   interval,
		batchSize:  defaultBatchSize,
		logger:     logger,
	}
}

// Run starts the polling loop. Blocks until ctx is cancelled.
func (c *Consumer) Run(ctx context.Context) {
	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			c.poll(ctx)
		}
	}
}

func (c *Consumer) poll(ctx context.Context) {
	primitives, err := c.outboxRepo.PullPending(ctx, c.batchSize)
	if err != nil {
		c.consecutiveErrors++
		if c.consecutiveErrors >= maxConsecutiveErrors {
			c.logger.Fatal("outbox consumer: too many consecutive pull failures, shutting down",
				shared_logger.Field{Key: "error", Value: err.Error()},
				shared_logger.Field{Key: "consecutive_errors", Value: c.consecutiveErrors},
			)
		}
		c.logger.Error("outbox consumer: pull pending failed",
			shared_logger.Field{Key: "error", Value: err.Error()},
		)
		return
	}
	c.consecutiveErrors = 0
	if len(primitives) == 0 {
		return
	}

	events := make([]event2.Event, 0, len(primitives))
	for _, p := range primitives {
		evt, err := c.registry.Build(p)
		if err != nil {
			c.logger.Error("outbox consumer: failed to build event",
				shared_logger.Field{Key: "event_type", Value: p.Type},
				shared_logger.Field{Key: "event_id", Value: p.ID},
				shared_logger.Field{Key: "error", Value: err.Error()},
			)
			continue
		}
		events = append(events, evt)
	}

	if len(events) > 0 {
		c.inner.Publish(ctx, events)
	}
}
