package persistance

import (
	"context"
	"github.com/isaac-olle/gas/internal/shared/domain/aggregate"
	"github.com/isaac-olle/gas/internal/shared/domain/criteria"
	"github.com/isaac-olle/gas/internal/shared/domain/vo/ids"
)

type Repository[T aggregate.Aggregate, I ids.Id] interface {
	Save(ctx context.Context, aggregate T) error
	SearchById(ctx context.Context, id I) (T, error)
	Matching(ctx context.Context, criteria criteria.Criteria) ([]T, error)
	SoftDelete(ctx context.Context, id I) error
	HardDelete(ctx context.Context, id I) error
}
