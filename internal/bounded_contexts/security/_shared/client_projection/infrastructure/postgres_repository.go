package infrastructure

import (
	"context"
	"errors"
	"fmt"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/_shared/client_projection"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"github.com/isaac-olle/gas-backoffice/internal/shared/persistance/infrastructure"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ClientProjectionRepository struct {
	postgres *infrastructure.PostgresRepository[*client_projection.ClientProjection]
}

func NewClientProjectionRepository(db *pgxpool.Pool) *ClientProjectionRepository {
	return &ClientProjectionRepository{
		postgres: infrastructure.NewPostgresRepository[*client_projection.ClientProjection](db),
	}
}

func (r *ClientProjectionRepository) Find(ctx context.Context, id ids.ClientId) (*client_projection.ClientProjection, error) {
	query := `SELECT id, subscription_id FROM backoffice.clients WHERE id = $1 AND deleted_at IS NULL`
	return r.postgres.Find(ctx, query, r.fromRow, id.String())
}

func (r *ClientProjectionRepository) fromRow(row pgx.Row) (*client_projection.ClientProjection, error) {
	var p client_projection.ClientProjection
	err := row.Scan(&p.ClientId, &p.SubscriptionId)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("this client projection does not exist")
	}
	if err != nil {
		return nil, fmt.Errorf("error scanning client projection: %w", err)
	}
	return &p, nil
}
