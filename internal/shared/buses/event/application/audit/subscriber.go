package audit

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/isaac-olle/gas/internal/shared/buses/event"
	"github.com/isaac-olle/gas/internal/shared/domain/vo"
)

type Subscriber struct {
	repo event.AuditRepository
}

func NewSubscriber(repo event.AuditRepository) *Subscriber {
	return &Subscriber{repo: repo}
}

func (s *Subscriber) SubscribedTo() event.Type {
	return event.AllTypes
}

func (s *Subscriber) On(ctx context.Context, evt event.Event) error {
	p := evt.Primitives()

	payload, err := json.Marshal(p)
	if err != nil {
		return fmt.Errorf("audit: marshal event payload: %w", err)
	}

	rec := event.AuditRecord{
		EventID:       p.ID,
		EventType:     p.Type,
		AggregateID:   p.AggregateID,
		AggregateType: p.AggregateType,
		Payload:       payload,
		OccurredAt:    p.OccurredOn,
	}

	if actor, ok := vo.ActorFromContext[vo.Actor](ctx); ok {
		prims := actor.ToPrimitives()
		rec.ActorID = &prims.Id
		rec.ActorRole = &prims.Role
	}

	return s.repo.Save(ctx, rec)
}
