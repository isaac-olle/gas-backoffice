package infrastructure

import (
	"context"
	"github.com/isaac-olle/gas/internal/shared/domain/aggregate"
	"github.com/isaac-olle/gas/internal/shared/domain/criteria"
	"github.com/isaac-olle/gas/internal/shared/domain/vo/ids"
	"github.com/isaac-olle/gas/internal/shared/persistance"
)

// InMemoryContextDecorator wraps a Repository and caches SearchById results
// in the request context so the same aggregate is never fetched twice per request.
//
// Usage: call WithCache(ctx) in a middleware to initialise the cache, then pass
// the enriched ctx down to handlers. Any SearchById call on a context without a
// cache falls through to the inner repository transparently.
type InMemoryContextDecorator[T aggregate.Aggregate, I ids.Id] struct {
	inner persistance.Repository[T, I]
}

func NewInMemoryContextDecorator[T aggregate.Aggregate, I ids.Id](
	inner persistance.Repository[T, I],
) *InMemoryContextDecorator[T, I] {
	return &InMemoryContextDecorator[T, I]{inner: inner}
}

// WithCache returns a new context with an empty cache bound to this decorator instance.
// Each decorator instance uses its own pointer as the context key, so multiple
// decorators in the same request never share a cache.
func (d *InMemoryContextDecorator[T, I]) WithCache(ctx context.Context) context.Context {
	return context.WithValue(ctx, d, make(map[string]any))
}

func (d *InMemoryContextDecorator[T, I]) SearchById(ctx context.Context, id I) (T, error) {
	if cache, ok := ctx.Value(d).(map[string]any); ok {
		if v, found := cache[id.String()]; found {
			return v.(T), nil
		}
		result, err := d.inner.SearchById(ctx, id)
		if err == nil {
			cache[id.String()] = result
		}
		return result, err
	}
	return d.inner.SearchById(ctx, id)
}

func (d *InMemoryContextDecorator[T, I]) Save(ctx context.Context, agg T) error {
	return d.inner.Save(ctx, agg)
}

func (d *InMemoryContextDecorator[T, I]) Matching(ctx context.Context, c criteria.Criteria) ([]T, error) {
	return d.inner.Matching(ctx, c)
}

func (d *InMemoryContextDecorator[T, I]) SoftDelete(ctx context.Context, id I) error {
	return d.inner.SoftDelete(ctx, id)
}

func (d *InMemoryContextDecorator[T, I]) HardDelete(ctx context.Context, id I) error {
	return d.inner.HardDelete(ctx, id)
}
