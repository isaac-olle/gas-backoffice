package event

import "context"

type Bus interface {
	Publish(context.Context, []Event)
}

type AuditRepository interface {
	Save(ctx context.Context, r AuditRecord) error
}

type OutboxRepository interface {
	Save(ctx context.Context, events []Event) error
	PullPending(ctx context.Context, limit int) ([]*Primitives, error)
}
