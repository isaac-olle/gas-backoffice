package vo

import (
	"context"
	"strings"
)

type scopesContextKey struct{}

func WithScopes(ctx context.Context, scopes []string) context.Context {
	return context.WithValue(ctx, scopesContextKey{}, scopes)
}

func ScopesFromContext(ctx context.Context) ([]string, bool) {
	scopes, ok := ctx.Value(scopesContextKey{}).([]string)
	return scopes, ok
}

func HasScope(ctx context.Context, scope string) bool {
	scopes, ok := ScopesFromContext(ctx)
	if !ok {
		return false
	}
	for _, s := range scopes {
		if s == scope {
			return true
		}
	}
	return false
}

func ParseScopes(raw string) []string {
	if raw == "" {
		return nil
	}
	return strings.Fields(raw)
}
