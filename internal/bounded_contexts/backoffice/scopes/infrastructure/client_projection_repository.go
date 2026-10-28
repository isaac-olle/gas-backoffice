package infrastructure

import (
	"context"
	"errors"
	"fmt"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/scopes"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	sharedInfrastructure "github.com/isaac-olle/gas-backoffice/internal/shared/persistance/infrastructure"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ClientProjectionRepository struct {
	postgres *sharedInfrastructure.PostgresRepository[*scopes.ClientProjection]
}

func NewClientProjectionRepository(db *pgxpool.Pool) *ClientProjectionRepository {
	return &ClientProjectionRepository{
		postgres: sharedInfrastructure.NewPostgresRepository[*scopes.ClientProjection](db),
	}
}

func (r *ClientProjectionRepository) FindClientProjection(ctx context.Context, id ids.ClientId) (*scopes.ClientProjection, error) {
	return r.postgres.Find(ctx,
		`SELECT id, subscription_id FROM backoffice.clients WHERE id = $1 AND deleted_at IS NULL`,
		r.fromRow, id.String())
}

func (r *ClientProjectionRepository) fromRow(row pgx.Row) (*scopes.ClientProjection, error) {
	var rawClientID, rawSubscriptionID string
	if err := row.Scan(&rawClientID, &rawSubscriptionID); errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("error scanning client projection: %w", err)
	}
	clientID, err := ids.NewClientIdFromString(rawClientID)
	if err != nil {
		return nil, fmt.Errorf("invalid client id in projection: %w", err)
	}
	subscriptionID, err := ids.NewSubscriptionIdFromString(rawSubscriptionID)
	if err != nil {
		return nil, fmt.Errorf("invalid subscription id in projection: %w", err)
	}
	return &scopes.ClientProjection{ClientId: clientID, SubscriptionId: subscriptionID}, nil
}
