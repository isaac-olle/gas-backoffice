package infrastructure

import (
	"context"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresClientRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresClientRepository(pool *pgxpool.Pool) *PostgresClientRepository {
	return &PostgresClientRepository{pool: pool}
}

func (r *PostgresClientRepository) FindById(ctx context.Context, clientID string) (*session.Client, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT
			c.id,
			c.name,
			c.tenant_id,
			c.admin_id,
			c.allowed_flows,
			c.redirect_uris,
			s.id,
			s.name,
			s.protected_resource_id,
			pr.name
		FROM clients c
		LEFT JOIN scopes s
			ON s.id = ANY(c.allowed_scopes) AND s.deleted_at IS NULL
		LEFT JOIN protected_resources pr
			ON pr.id = s.protected_resource_id AND pr.deleted_at IS NULL
		WHERE c.id = $1 AND c.deleted_at IS NULL
	`, clientID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var client *session.Client

	for rows.Next() {
		var (
			id           string
			name         string
			tenantID     string
			adminID      string
			allowedFlows []string
			redirectURIs []string
			scopeID      *string
			scopeName    *string
			prID         *string
			prName       *string
		)

		if err = rows.Scan(&id, &name, &tenantID, &adminID, &allowedFlows, &redirectURIs, &scopeID, &scopeName, &prID, &prName); err != nil {
			return nil, err
		}

		if client == nil {
			client = session.NewClient(id, name, tenantID, adminID, allowedFlows, redirectURIs, nil)
		}

		if scopeID != nil {
			client.Scopes = append(client.Scopes, session.NewScope(*scopeID, *scopeName, *prID, *prName))
		}
	}

	return client, rows.Err()
}
