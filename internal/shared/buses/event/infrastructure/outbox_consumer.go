package infrastructure

import (
	"context"
	event2 "github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/authors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	shared_logger "github.com/isaac-olle/gas-backoffice/internal/shared/logger"
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
			c.Poll(ctx)
		}
	}
}

// Poll pulls one batch of pending outbox events and dispatches each to
// subscribers. Exported so it can be driven directly in tests without waiting
// on the ticker.
func (c *Consumer) Poll(ctx context.Context) {
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
		c.inner.Publish(c.withAuthor(ctx, p), []event2.Event{evt})
	}
}

// withAuthor rebuilds the vo.Author that published this event (persisted
// alongside it in the outbox row) and injects it into ctx, so handlers that
// write to the database still get the right author_id/author_role. Events
// published without an author (e.g. system-generated) pass ctx through as-is.
func (c *Consumer) withAuthor(ctx context.Context, p *event2.Primitives) context.Context {
	if p.AuthorId == "" {
		return ctx
	}
	role, err := vo.NewRole(p.AuthorRole)
	if err != nil {
		c.logger.Error("outbox consumer: invalid author role",
			shared_logger.Field{Key: "event_id", Value: p.ID},
			shared_logger.Field{Key: "author_role", Value: p.AuthorRole},
			shared_logger.Field{Key: "error", Value: err.Error()},
		)
		return ctx
	}
	// The outbox only persists author_id/author_role (all subscribers need for
	// auditing), so the rebuilt author carries no AdminId.
	authorId, err := ids.NewUserIdFromString(p.AuthorId)
	if err != nil {
		c.logger.Error("outbox consumer: invalid author id",
			shared_logger.Field{Key: "event_id", Value: p.ID},
			shared_logger.Field{Key: "author_id", Value: p.AuthorId},
			shared_logger.Field{Key: "error", Value: err.Error()},
		)
		return ctx
	}
	return vo.WithAuthor(ctx, &authors.UserAuthor{Id: authorId, Role: role})
}
