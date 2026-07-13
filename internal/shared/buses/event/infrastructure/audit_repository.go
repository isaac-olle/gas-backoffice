package infrastructure

import (
	"context"
	"fmt"

	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/event"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresAuditRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresAuditRepository(pool *pgxpool.Pool) *PostgresAuditRepository {
	return &PostgresAuditRepository{pool: pool}
}

func (r *PostgresAuditRepository) Save(ctx context.Context, rec event.AuditRecord) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO audit_events (event_id, event_type, aggregate_id, aggregate_type, author_id, author_role, payload, occurred_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		 ON CONFLICT (event_id) DO NOTHING`,
		rec.EventID, rec.EventType, rec.AggregateID, rec.AggregateType,
		rec.AuthorID, rec.AuthorRole, rec.Payload, rec.OccurredAt,
	)
	if err != nil {
		return fmt.Errorf("audit: insert event %s: %w", rec.EventID, err)
	}
	return nil
}
