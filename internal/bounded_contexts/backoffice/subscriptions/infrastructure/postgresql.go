package infrastructure

import (
	"context"
	"errors"
	"fmt"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/subscriptions"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"github.com/isaac-olle/gas-backoffice/internal/shared/persistance"
	pginfra "github.com/isaac-olle/gas-backoffice/internal/shared/persistance/infrastructure"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct {
	postgres *pginfra.PostgresRepository[*subscriptions.Subscription]
}

func NewPostgresRepository(db *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{postgres: pginfra.NewPostgresRepository[*subscriptions.Subscription](db)}
}

func (r *PostgresRepository) Save(ctx context.Context, aggregate *subscriptions.Subscription) error {
	// Join the ambient transaction (set by the UnitOfWork around MultipleExecution)
	// when one exists, instead of always opening a new one: a brand-new
	// transaction can't see rows inserted but not yet committed by earlier
	// commands in the same batch (e.g. the admin row a subscription's
	// owner_id FK points at).
	if tx, ok := pginfra.TxFromContext(ctx); ok {
		return r.save(ctx, tx, aggregate)
	}

	return r.postgres.WithTx(ctx, func(tx pgx.Tx) error { return r.save(ctx, tx, aggregate) })
}

func (r *PostgresRepository) save(ctx context.Context, q pginfra.Querier, aggregate *subscriptions.Subscription) error {
	p := aggregate.ToPrimitives()
	authorId, authorRole := "", ""
	if author, ok := vo.AuthorFromContext(ctx); ok {
		a := author.ToPrimitives()
		authorId, authorRole = a.Id, a.Role
	}
	err := q.QueryRow(ctx, `INSERT INTO backoffice.subscriptions (id,name,owner_id,is_system,version,author_id,author_role)
		VALUES ($1,$2,$3,$4,$5,$6,$7) ON CONFLICT (id) DO UPDATE SET name=EXCLUDED.name, is_system=EXCLUDED.is_system, version=EXCLUDED.version+1
		WHERE backoffice.subscriptions.version=EXCLUDED.version RETURNING version`, p.Id, p.Name, p.OwnerId, p.IsSystem, p.Version, authorId, authorRole).Scan(new(int))
	if errors.Is(err, pgx.ErrNoRows) {
		return persistance.ErrConcurrentModification
	}
	if err != nil {
		return fmt.Errorf("save subscription: %w", err)
	}
	if _, err = q.Exec(ctx, `DELETE FROM backoffice.subscription_members WHERE subscription_id=$1`, p.Id); err != nil {
		return fmt.Errorf("replace subscription members: %w", err)
	}
	for _, member := range p.Members {
		if _, err = q.Exec(ctx, `INSERT INTO backoffice.subscription_members (subscription_id,user_id,subscription_role) VALUES ($1,$2,$3)`, p.Id, member.UserId, member.SubscriptionRole); err != nil {
			return fmt.Errorf("save subscription member: %w", err)
		}
	}
	return nil
}

func (r *PostgresRepository) SearchById(ctx context.Context, id ids.SubscriptionId) (*subscriptions.Subscription, error) {
	var p subscriptions.SubscriptionPrimitives
	var memberIds, memberRoles []string
	err := r.postgres.QueryRow(ctx, `SELECT s.id,s.name,s.owner_id,s.is_system,
		COALESCE(array_agg(m.user_id) FILTER (WHERE m.user_id IS NOT NULL),'{}'),
		COALESCE(array_agg(m.subscription_role) FILTER (WHERE m.user_id IS NOT NULL),'{}'),s.version
		FROM backoffice.subscriptions s LEFT JOIN backoffice.subscription_members m ON m.subscription_id=s.id
		WHERE s.id=$1 AND s.deleted_at IS NULL GROUP BY s.id`, id.String()).Scan(&p.Id, &p.Name, &p.OwnerId, &p.IsSystem, &memberIds, &memberRoles, &p.Version)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("find subscription: %w", err)
	}
	p.Members = make([]subscriptions.MembershipPrimitives, 0, len(memberIds))
	for i, userId := range memberIds {
		p.Members = append(p.Members, subscriptions.MembershipPrimitives{UserId: userId, SubscriptionRole: memberRoles[i]})
	}
	return subscriptions.FromPrimitives(p)
}

func (r *PostgresRepository) List(ctx context.Context, subscriptionIds []ids.SubscriptionId, all bool) ([]subscriptions.Summary, error) {
	if !all && len(subscriptionIds) == 0 {
		return []subscriptions.Summary{}, nil
	}

	query := `SELECT id, name, is_system FROM backoffice.subscriptions WHERE deleted_at IS NULL`
	args := make([]any, 0, 1)
	if !all {
		rawIds := make([]string, 0, len(subscriptionIds))
		for _, id := range subscriptionIds {
			rawIds = append(rawIds, id.String())
		}
		query += ` AND id = ANY($1)`
		args = append(args, rawIds)
	}
	query += ` ORDER BY name, id`

	rows, err := r.postgres.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list subscriptions: %w", err)
	}
	defer rows.Close()

	result := make([]subscriptions.Summary, 0)
	for rows.Next() {
		var item subscriptions.Summary
		if err = rows.Scan(&item.Id, &item.Name, &item.IsSystem); err != nil {
			return nil, fmt.Errorf("scan subscription summary: %w", err)
		}
		result = append(result, item)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate subscriptions: %w", err)
	}
	return result, nil
}

var _ subscriptions.Repository = (*PostgresRepository)(nil)
