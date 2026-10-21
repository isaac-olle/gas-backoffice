package infrastructure

import (
	"context"
	"errors"
	"fmt"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/subscriptions"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"github.com/isaac-olle/gas-backoffice/internal/shared/persistance"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct{ db *pgxpool.Pool }

func NewPostgresRepository(db *pgxpool.Pool) *PostgresRepository { return &PostgresRepository{db: db} }

func (r *PostgresRepository) Save(ctx context.Context, aggregate *subscriptions.Subscription) error {
	p := aggregate.ToPrimitives()
	authorId, authorRole := "", ""
	if author, ok := vo.AuthorFromContext[vo.Author](ctx); ok {
		a := author.ToPrimitives()
		authorId, authorRole = a.Id, a.Role
	}
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin subscription transaction: %w", err)
	}
	defer tx.Rollback(ctx)
	err = tx.QueryRow(ctx, `INSERT INTO backoffice.subscriptions (id,name,created_by,version,author_id,author_role)
		VALUES ($1,$2,$3,$4,$5,$6) ON CONFLICT (id) DO UPDATE SET name=EXCLUDED.name, version=EXCLUDED.version+1
		WHERE backoffice.subscriptions.version=EXCLUDED.version RETURNING version`, p.Id, p.Name, p.CreatedBy, p.Version, authorId, authorRole).Scan(new(int))
	if errors.Is(err, pgx.ErrNoRows) {
		return persistance.ErrConcurrentModification
	}
	if err != nil {
		return fmt.Errorf("save subscription: %w", err)
	}
	if _, err = tx.Exec(ctx, `DELETE FROM backoffice.subscription_members WHERE subscription_id=$1`, p.Id); err != nil {
		return fmt.Errorf("replace subscription members: %w", err)
	}
	for _, userId := range p.MemberIds {
		if _, err = tx.Exec(ctx, `INSERT INTO backoffice.subscription_members (subscription_id,user_id) VALUES ($1,$2)`, p.Id, userId); err != nil {
			return fmt.Errorf("save subscription member: %w", err)
		}
	}
	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit subscription: %w", err)
	}
	return nil
}

func (r *PostgresRepository) SearchById(ctx context.Context, id ids.SubscriptionId) (*subscriptions.Subscription, error) {
	var p subscriptions.SubscriptionPrimitives
	err := r.db.QueryRow(ctx, `SELECT s.id,s.name,s.created_by,
		COALESCE(array_agg(m.user_id) FILTER (WHERE m.user_id IS NOT NULL),'{}'),s.version
		FROM backoffice.subscriptions s LEFT JOIN backoffice.subscription_members m ON m.subscription_id=s.id
		WHERE s.id=$1 AND s.deleted_at IS NULL GROUP BY s.id`, id.String()).Scan(&p.Id, &p.Name, &p.CreatedBy, &p.MemberIds, &p.Version)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("find subscription: %w", err)
	}
	return subscriptions.FromPrimitives(p)
}

var _ subscriptions.Repository = (*PostgresRepository)(nil)
