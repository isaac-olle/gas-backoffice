package persistance

import (
	"context"
	"fmt"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresClientCredentialsSessionRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresClientCredentialsSessionRepository(pool *pgxpool.Pool) *PostgresClientCredentialsSessionRepository {
	return &PostgresClientCredentialsSessionRepository{pool: pool}
}

func (r *PostgresClientCredentialsSessionRepository) Save(ctx context.Context, sess *session.ClientCredentialsSession) error {
	p := sess.ToPrimitives()

	_, err := r.pool.Exec(ctx, `
		INSERT INTO protocol_client_credentials_sessions (
			id, client_id, tenant_id, admin_id, scope, issued_at
		) VALUES ($1,$2,$3,$4,$5,$6)
	`, p.ID, p.ClientID, p.TenantID, p.AdminID, p.Scope, p.IssuedAt)
	if err != nil {
		return fmt.Errorf("exec insert client credentials session: %w", err)
	}
	return nil
}
