package infrastructure

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/auth/_shared/user_projection"
	"github.com/isaac-olle/gas-backoffice/internal/shared/persistance/infrastructure"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserProjectionRepository struct {
	postgres *infrastructure.PostgresRepository[*user_projection.UserProjection]
}

func NewUserProjectionRepository(db *pgxpool.Pool) *UserProjectionRepository {
	return &UserProjectionRepository{
		postgres: infrastructure.NewPostgresRepository[*user_projection.UserProjection](db),
	}
}

func (r *UserProjectionRepository) FindByEmail(ctx context.Context, email string) (*user_projection.UserProjection, error) {
	query := `
		SELECT u.id, u.email, u.system_role, CASE WHEN u.system_role = 'system_administrator' THEN u.id ELSE u.author_id END AS admin_id, u.subscription_ids, u.memberships
		FROM backoffice.users u
		WHERE u.email = $1 AND u.deleted_at IS NULL
	`
	return r.postgres.Find(ctx, query, r.fromRow, email)
}

func (r *UserProjectionRepository) fromRow(row pgx.Row) (*user_projection.UserProjection, error) {
	var p user_projection.UserProjection
	var memberships []byte
	err := row.Scan(&p.UserId, &p.Email, &p.SystemRole, &p.AdminId, &p.SubscriptionIds, &memberships)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("no user found for that email")
	}
	if err != nil {
		return nil, fmt.Errorf("error scanning user projection: %w", err)
	}
	if err = json.Unmarshal(memberships, &p.Memberships); err != nil {
		return nil, fmt.Errorf("unmarshal memberships: %w", err)
	}
	return &p, nil
}
