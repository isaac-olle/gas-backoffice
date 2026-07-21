package bootstrap

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

func resolveApplicationIDs(ctx context.Context, pool *pgxpool.Pool, backofficePRName string) (backofficePrId, openIDScopeId string, err error) {
	err = pool.QueryRow(ctx,
		`SELECT id FROM backoffice.protected_resources WHERE name = $1 AND is_system = TRUE AND deleted_at IS NULL LIMIT 1`,
		backofficePRName,
	).Scan(&backofficePrId)
	if err != nil {
		return "", "", fmt.Errorf("resolve backoffice protected resource id: %w", err)
	}

	err = pool.QueryRow(ctx,
		`SELECT id FROM backoffice.scopes WHERE name = 'openid' AND protected_resource_id = $1 AND deleted_at IS NULL LIMIT 1`,
		backofficePrId,
	).Scan(&openIDScopeId)
	if err != nil {
		return "", "", fmt.Errorf("resolve openid scope id: %w", err)
	}

	return backofficePrId, openIDScopeId, nil
}
