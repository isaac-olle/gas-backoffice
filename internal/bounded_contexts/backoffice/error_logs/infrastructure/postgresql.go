package infrastructure

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/error_logs"
	"github.com/isaac-olle/gas-backoffice/internal/shared/persistance/infrastructure"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresErrorLogRepository struct {
	postgres *infrastructure.PostgresRepository[error_logs.ErrorLog]
}

func NewPostgresErrorLogRepository(db *pgxpool.Pool) *PostgresErrorLogRepository {
	return &PostgresErrorLogRepository{
		postgres: infrastructure.NewPostgresRepository[error_logs.ErrorLog](db),
	}
}

func (r *PostgresErrorLogRepository) Save(ctx context.Context, errorLog error_logs.ErrorLog) error {
	query := `
		INSERT INTO public.errors (error_id, type, code, description, details, cause, stack_trace, request_id, occurred_at,
			author_id, author_role)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NULLIF($10, ''), NULLIF($11, ''))
		ON CONFLICT (error_id) DO NOTHING
	`
	details, err := json.Marshal(errorLog.Details)
	if err != nil {
		return fmt.Errorf("marshal error log details: %w", err)
	}
	stackTrace, err := json.Marshal(errorLog.StackTrace)
	if err != nil {
		return fmt.Errorf("marshal error log stack trace: %w", err)
	}

	return r.postgres.Save(ctx, query,
		errorLog.ErrorId, errorLog.Type, errorLog.Code, errorLog.Description,
		details, errorLog.Cause, stackTrace, errorLog.RequestId, errorLog.OccurredAt,
		errorLog.AuthorId, errorLog.AuthorRole,
	)
}

func (r *PostgresErrorLogRepository) Latest(ctx context.Context, limit int) ([]error_logs.ErrorLog, error) {
	query := `SELECT error_id, type, code, description, details, request_id, occurred_at,
		COALESCE(author_id, ''), COALESCE(author_role, '')
		FROM public.errors
		ORDER BY occurred_at DESC NULLS LAST
		LIMIT $1`
	return r.postgres.FindMany(ctx, query, r.fromRow, limit)
}

func (r *PostgresErrorLogRepository) fromRow(row pgx.Row) (error_logs.ErrorLog, error) {
	var e error_logs.ErrorLog
	var requestId *string
	if err := row.Scan(&e.ErrorId, &e.Type, &e.Code, &e.Description, &e.Details, &requestId, &e.OccurredAt,
		&e.AuthorId, &e.AuthorRole); err != nil {
		return error_logs.ErrorLog{}, fmt.Errorf("error scanning error log: %w", err)
	}
	if requestId != nil {
		e.RequestId = *requestId
	}
	return e, nil
}
