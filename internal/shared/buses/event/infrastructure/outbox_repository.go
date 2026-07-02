package infrastructure

import (
	"context"
	"encoding/json"
	"fmt"
	event2 "github.com/isaac-olle/gas/internal/shared/buses/event"
	"github.com/isaac-olle/gas/internal/shared/domain/vo"
	persistance "github.com/isaac-olle/gas/internal/shared/persistance/infrastructure"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository persists domain events into the outbox_event_bus table and
// exposes them for the poller to consume.
type Repository struct {
	postgres *persistance.PostgresRepository[event2.Event]
}

func NewOutboxRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{
		postgres: persistance.NewPostgresRepository[event2.Event](pool),
	}
}

// Save inserts domain events into the outbox.
// Participates in the ambient transaction from ctx if present.
// The actor is extracted from ctx if available.
func (r *Repository) Save(ctx context.Context, events []event2.Event) error {
	var actorId, actorRole *string
	if actor, ok := vo.ActorFromContext[vo.Actor](ctx); ok {
		p := actor.ToPrimitives()
		actorId = &p.Id
		actorRole = &p.Role
	}

	for _, evt := range events {
		p := evt.Primitives()
		payload, err := json.Marshal(p)
		if err != nil {
			return fmt.Errorf("marshal event payload: %w", err)
		}
		if err = r.postgres.Save(ctx,
			`INSERT INTO outbox_event_bus (event_id, event_name, aggregate_id, payload, occurred_at, actor_id, actor_role)
			 VALUES ($1, $2, $3, $4, $5, $6, $7)`,
			p.ID, p.Type, p.AggregateID, payload, p.OccurredOn, actorId, actorRole,
		); err != nil {
			return fmt.Errorf("insert outbox event %s: %w", p.ID, err)
		}
	}
	return nil
}

// PullPending returns up to limit unprocessed events ordered by insertion time
// and deletes them atomically. Uses FOR UPDATE SKIP LOCKED so concurrent
// pollers never pick the same rows.
func (r *Repository) PullPending(ctx context.Context, limit int) ([]*event2.Primitives, error) {
	tx, err := r.postgres.DB().Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	rows, err := tx.Query(ctx,
		`SELECT event_id, payload
		 FROM outbox_event_bus
		 ORDER BY inserted_at ASC
		 LIMIT $1
		 FOR UPDATE SKIP LOCKED`,
		limit,
	)
	if err != nil {
		return nil, fmt.Errorf("query outbox: %w", err)
	}

	var primitives []*event2.Primitives
	var ids []string
	for rows.Next() {
		p, err := r.fromSqlToDomain(rows)
		if err != nil {
			rows.Close()
			return nil, err
		}
		primitives = append(primitives, p)
		ids = append(ids, p.ID)
	}
	rows.Close()

	if len(ids) == 0 {
		return nil, nil
	}

	if err = r.delete(ctx, tx, ids); err != nil {
		return nil, err
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit tx: %w", err)
	}

	return primitives, nil
}

func (r *Repository) fromSqlToDomain(row pgx.Row) (*event2.Primitives, error) {
	var id string
	var payload []byte
	if err := row.Scan(&id, &payload); err != nil {
		return nil, fmt.Errorf("scan outbox event: %w", err)
	}
	var p event2.Primitives
	if err := json.Unmarshal(payload, &p); err != nil {
		return nil, fmt.Errorf("unmarshal outbox event: %w", err)
	}
	return &p, nil
}

func (r *Repository) delete(ctx context.Context, tx pgx.Tx, ids []string) error {
	if _, err := tx.Exec(ctx,
		`DELETE FROM outbox_event_bus WHERE event_id = ANY($1)`,
		ids,
	); err != nil {
		return fmt.Errorf("delete outbox events: %w", err)
	}
	return nil
}
