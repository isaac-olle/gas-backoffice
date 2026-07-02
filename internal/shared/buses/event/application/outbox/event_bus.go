package outbox

import (
	"context"

	"github.com/isaac-olle/gas/internal/shared/buses/event"
	shared_logger "github.com/isaac-olle/gas/internal/shared/logger"
)

type EventBus struct {
	outboxRepo event.OutboxRepository
	logger     shared_logger.Logger
}

func NewEventBus(outboxRepo event.OutboxRepository, logger shared_logger.Logger) *EventBus {
	return &EventBus{
		outboxRepo: outboxRepo,
		logger:     logger,
	}
}

func (b *EventBus) Publish(ctx context.Context, events []event.Event) {
	if len(events) == 0 {
		return
	}
	if err := b.outboxRepo.Save(ctx, events); err != nil {
		b.logger.Error("outbox: failed to save events",
			shared_logger.Field{Key: "error", Value: err.Error()},
		)
	}
}
