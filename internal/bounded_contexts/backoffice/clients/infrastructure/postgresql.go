package infrastructure

import (
	"context"
	"errors"
	"fmt"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/criteria"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"github.com/isaac-olle/gas-backoffice/internal/shared/persistance/infrastructure"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresClientRepository struct {
	postgres *infrastructure.PostgresRepository[*clients.Client]
}

func NewPostgresClientRepository(db *pgxpool.Pool) *PostgresClientRepository {
	return &PostgresClientRepository{
		postgres: infrastructure.NewPostgresRepository[*clients.Client](db),
	}
}

func (r *PostgresClientRepository) Save(ctx context.Context, client *clients.Client) error {
	p := client.ToPrimitives()
	authorId, authorRole := actorPrimitives(ctx)
	query := `
		INSERT INTO backoffice.clients (id, name, admin_id, user_id, allowed_flows, allowed_scopes, redirect_uris, key_ids, secret_ids, author_id, author_role)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		ON CONFLICT (id) DO UPDATE SET
			name = EXCLUDED.name,
			allowed_flows = EXCLUDED.allowed_flows,
			allowed_scopes = EXCLUDED.allowed_scopes,
			redirect_uris = EXCLUDED.redirect_uris,
			key_ids = EXCLUDED.key_ids,
			secret_ids = EXCLUDED.secret_ids
	`
	return r.postgres.Save(ctx, query, p.Id, p.Name, p.AdminId, p.UserId, p.AllowedFlows, p.AllowedScopes, p.RedirectUris, p.KeyIds, p.SecretIds, authorId, authorRole)
}

func actorPrimitives(ctx context.Context) (string, string) {
	actor, ok := vo.ActorFromContext[vo.Actor](ctx)
	if !ok {
		return "", ""
	}
	p := actor.ToPrimitives()
	return p.Id, p.Role
}

func (r *PostgresClientRepository) SearchById(ctx context.Context, id ids.ClientId) (*clients.Client, error) {
	query := `SELECT id, name, admin_id, user_id, allowed_flows, allowed_scopes, redirect_uris, key_ids, secret_ids FROM backoffice.clients WHERE id = $1 AND deleted_at IS NULL`
	return r.postgres.Find(ctx, query, r.fromRow, id.String())
}

func (r *PostgresClientRepository) Matching(ctx context.Context, _ criteria.Criteria) ([]*clients.Client, error) {
	return nil, fmt.Errorf("not implemented")
}

func (r *PostgresClientRepository) SoftDelete(ctx context.Context, id ids.ClientId) error {
	query := `UPDATE backoffice.clients SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL`
	return r.postgres.Delete(ctx, query, id.String())
}

func (r *PostgresClientRepository) HardDelete(ctx context.Context, id ids.ClientId) error {
	query := `DELETE FROM backoffice.clients WHERE id = $1`
	return r.postgres.Delete(ctx, query, id.String())
}

func (r *PostgresClientRepository) fromRow(row pgx.Row) (*clients.Client, error) {
	var p clients.ClientPrimitives
	err := row.Scan(&p.Id, &p.Name, &p.AdminId, &p.UserId, &p.AllowedFlows, &p.AllowedScopes, &p.RedirectUris, &p.KeyIds, &p.SecretIds)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error scanning client: %w", err)
	}
	return clients.FromPrimitives(p)
}
