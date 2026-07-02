package vo

import "context"

type Ownership interface{}

type ownershipContextKey struct{}

func WithOwnership(ctx context.Context, ownership Ownership) context.Context {
	return context.WithValue(ctx, ownershipContextKey{}, ownership)
}

func OwnershipFromContext[O Ownership](ctx context.Context) (O, bool) {
	ownership, ok := ctx.Value(ownershipContextKey{}).(O)
	return ownership, ok
}
