package vo

import (
	"context"
)

type authorContextKey struct{}

type Author interface {
	EqualsRole(role Role) bool
	ToPrimitives() *AuthorPrimitives
}

func WithAuthor(ctx context.Context, author Author) context.Context {
	return context.WithValue(ctx, authorContextKey{}, author)
}

func AuthorFromContext[A Author](ctx context.Context) (A, bool) {
	author, ok := ctx.Value(authorContextKey{}).(A)
	return author, ok
}

type AuthorPrimitives struct {
	Id   string
	Role string
}

type OwnedAuthor[O Ownership] interface {
	Author
	Ownership() O
}
