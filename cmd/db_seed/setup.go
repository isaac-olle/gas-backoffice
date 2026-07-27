package main

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo"
)

var systemScopes = []string{
	vo.ScopeUsersRead,
	vo.ScopeUsersReadAll,
	vo.ScopeUsersCreate,
	vo.ScopeUsersRename,
	vo.ScopeUsersDelete,
	vo.ScopeProtectedResourcesCreate,
	vo.ScopeProtectedResourcesRead,
	vo.ScopeProtectedResourcesRename,
	vo.ScopeProtectedResourcesDelete,
	vo.ScopeScopesCreate,
	vo.ScopeScopesRead,
	vo.ScopeScopesDelete,
	vo.ScopeScopesAddAllowedClient,
	vo.ScopeScopesRemoveAllowedClient,
	vo.ScopeClientsCreate,
	vo.ScopeClientsRead,
	vo.ScopeClientsDelete,
	vo.ScopeClientsAddFlow,
	vo.ScopeClientsRemoveFlow,
	vo.ScopeClientsAddRedirectUri,
	"openid",
}

func runSetup(ctx context.Context, pool *pgxpool.Pool, adminID, backofficePRName string) error {
	prID := uuid.New().String()
	keyPairID := uuid.New().String()

	_ = pool.QueryRow(ctx,
		`SELECT id FROM backoffice.protected_resources WHERE name = $1 AND is_system = TRUE AND deleted_at IS NULL LIMIT 1`,
		backofficePRName,
	).Scan(&prID)

	_ = pool.QueryRow(ctx,
		`SELECT id FROM security.key_pairs WHERE protected_resource_id = $1 LIMIT 1`,
		prID,
	).Scan(&keyPairID)

	type scopeEntry struct {
		id   string
		name string
	}
	scopes := make([]scopeEntry, len(systemScopes))
	for i, name := range systemScopes {
		scopes[i] = scopeEntry{id: uuid.New().String(), name: name}
		_ = pool.QueryRow(ctx,
			`SELECT id FROM backoffice.scopes WHERE name = $1 AND protected_resource_id = $2 AND deleted_at IS NULL LIMIT 1`,
			name, prID,
		).Scan(&scopes[i].id)
	}

	privKey, pubKeyPEM, err := generateRSAKeyPair()
	if err != nil {
		return fmt.Errorf("generate key pair: %w", err)
	}
	privKeyPEM, err := marshalPrivateKey(privKey)
	if err != nil {
		return fmt.Errorf("marshal private key: %w", err)
	}

	tx, err := pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	allScopeIDs := make([]string, len(scopes))
	for i, s := range scopes {
		allScopeIDs[i] = s.id
	}

	if _, err = tx.Exec(ctx, `
		INSERT INTO backoffice.protected_resources (id, name, user_id, scope_ids, is_system, author_id, author_role)
		VALUES ($1, $2, $3, $4, TRUE, $5, 'admin')
		ON CONFLICT (id) DO UPDATE SET scope_ids = EXCLUDED.scope_ids, author_id = EXCLUDED.author_id, author_role = EXCLUDED.author_role`,
		prID, backofficePRName, adminID, allScopeIDs, adminID); err != nil {
		return fmt.Errorf("insert system protected resource: %w", err)
	}

	if _, err = tx.Exec(ctx, `
		INSERT INTO security.key_pairs (id, protected_resource_id, algorithm, public_key, private_key, author_id, author_role)
		VALUES ($1, $2, $3, $4, $5, $6, 'admin')
		ON CONFLICT (id) DO UPDATE SET author_id = EXCLUDED.author_id, author_role = EXCLUDED.author_role`,
		keyPairID, prID, "RS256", pubKeyPEM, privKeyPEM, adminID); err != nil {
		return fmt.Errorf("insert key pair: %w", err)
	}

	for _, s := range scopes {
		if _, err = tx.Exec(ctx, `
			INSERT INTO backoffice.scopes (id, name, protected_resource_id, user_id, is_system, author_id, author_role)
			VALUES ($1, $2, $3, $4, TRUE, $5, 'admin')
			ON CONFLICT (protected_resource_id, name) WHERE deleted_at IS NULL DO UPDATE SET author_id = EXCLUDED.author_id, author_role = EXCLUDED.author_role`,
			s.id, s.name, prID, adminID, adminID); err != nil {
			return fmt.Errorf("insert scope %q: %w", s.name, err)
		}
	}

	if _, err = tx.Exec(ctx, `
		INSERT INTO backoffice.users (id, admin_id, name, role, is_system, protected_resource_ids, author_id, author_role)
		VALUES ($1, $1, 'Gas Seed Admin', 'admin', TRUE, $2, $1, 'admin')
		ON CONFLICT (id) DO UPDATE SET protected_resource_ids = EXCLUDED.protected_resource_ids, author_id = EXCLUDED.author_id, author_role = EXCLUDED.author_role`,
		adminID, []string{prID}); err != nil {
		return fmt.Errorf("upsert admin user: %w", err)
	}

	return tx.Commit(ctx)
}
