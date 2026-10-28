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
			c.allowed_grant_types,
			c.redirect_uris,
			c.public,
			s.id,
			s.name,
			s.protected_resource_id,
			pr.name
		FROM backoffice.clients c
		LEFT JOIN backoffice.scopes s
			ON $1 = ANY(s.allowed_client_ids) AND s.deleted_at IS NULL AND s.enabled = TRUE
		LEFT JOIN backoffice.protected_resources pr
			ON pr.id = s.protected_resource_id AND pr.deleted_at IS NULL AND pr.enabled = TRUE
		WHERE c.id = $1 AND c.deleted_at IS NULL AND c.enabled = TRUE
	`, clientID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var client *session.Client

	for rows.Next() {
		var (
			id                string
			name              string
			allowedGrantTypes []string
			redirectURIs      []string
			public            bool
			scopeID           *string
			scopeName         *string
			prID              *string
			prName            *string
		)

		if err = rows.Scan(&id, &name, &allowedGrantTypes, &redirectURIs, &public, &scopeID, &scopeName, &prID, &prName); err != nil {
			return nil, err
		}

		if client == nil {
			client = session.NewClient(id, name, "", allowedGrantTypes, redirectURIs, nil, public)
		}

		// prName is nil when the scope's protected resource is missing or soft-deleted
		// (LEFT JOIN with deleted_at IS NULL): skip the scope rather than deref a nil pointer.
		if scopeID != nil && prName != nil {
			client.Scopes = append(client.Scopes, session.NewScope(*scopeID, *scopeName, *prID, *prName))
		}
	}

	return client, rows.Err()
}
