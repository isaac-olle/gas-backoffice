package infrastructure

import (
	"context"
	"fmt"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/event_logs"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresEventLogRepository struct{ pool *pgxpool.Pool }

func NewPostgresEventLogRepository(pool *pgxpool.Pool) *PostgresEventLogRepository {
	return &PostgresEventLogRepository{pool: pool}
}

func (r *PostgresEventLogRepository) Latest(ctx context.Context, limit int) ([]event_logs.EventLog, error) {
	rows, err := r.pool.Query(ctx, `SELECT event_id::text, event_type, aggregate_id, aggregate_type,
		author_id, author_role, payload, occurred_at, recorded_at
		FROM events.audit_events ORDER BY occurred_at DESC, event_id DESC LIMIT $1`, limit)
	if err != nil {
		return nil, fmt.Errorf("query event logs: %w", err)
	}
	defer rows.Close()
	result := make([]event_logs.EventLog, 0)
	for rows.Next() {
		var record event_logs.EventLog
		if err := rows.Scan(&record.EventId, &record.EventType, &record.AggregateId, &record.AggregateType,
			&record.AuthorId, &record.AuthorRole, &record.Payload, &record.OccurredAt, &record.RecordedAt); err != nil {
			return nil, fmt.Errorf("scan event log: %w", err)
		}
		result = append(result, record)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("read event logs: %w", err)
	}
	return result, nil
}
