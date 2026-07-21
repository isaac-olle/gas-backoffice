package ids

import "context"

type userIdContextKey struct{}
type clientIdContextKey struct{}
type userAssertionKeyIdContextKey struct{}
type clientAssertionKeyIdContextKey struct{}
type clientSecretIdContextKey struct{}
type protectedResourceIdContextKey struct{}
type scopeIdContextKey struct{}

func WithUserAssertionKeyId(ctx context.Context, id UserAssertionKeyId) context.Context {
	return context.WithValue(ctx, userAssertionKeyIdContextKey{}, id)
}

func UserAssertionKeyIdFromContext(ctx context.Context) (UserAssertionKeyId, bool) {
	id, ok := ctx.Value(userAssertionKeyIdContextKey{}).(UserAssertionKeyId)
	return id, ok
}

func WithUserId(ctx context.Context, id UserId) context.Context {
	return context.WithValue(ctx, userIdContextKey{}, id)
}

func UserIdFromContext(ctx context.Context) (UserId, bool) {
	id, ok := ctx.Value(userIdContextKey{}).(UserId)
	return id, ok
}

func WithClientAssertionKeyId(ctx context.Context, id ClientAssertionKeyId) context.Context {
	return context.WithValue(ctx, clientAssertionKeyIdContextKey{}, id)
}

func ClientAssertionKeyIdFromContext(ctx context.Context) (ClientAssertionKeyId, bool) {
	id, ok := ctx.Value(clientAssertionKeyIdContextKey{}).(ClientAssertionKeyId)
	return id, ok
}

func WithClientSecretId(ctx context.Context, id ClientSecretId) context.Context {
	return context.WithValue(ctx, clientSecretIdContextKey{}, id)
}

func ClientSecretIdFromContext(ctx context.Context) (ClientSecretId, bool) {
	id, ok := ctx.Value(clientSecretIdContextKey{}).(ClientSecretId)
	return id, ok
}

func WithClientId(ctx context.Context, id ClientId) context.Context {
	return context.WithValue(ctx, clientIdContextKey{}, id)
}

func ClientIdFromContext(ctx context.Context) (ClientId, bool) {
	id, ok := ctx.Value(clientIdContextKey{}).(ClientId)
	return id, ok
}

func WithProtectedResourceId(ctx context.Context, id ProtectedResourceId) context.Context {
	return context.WithValue(ctx, protectedResourceIdContextKey{}, id)
}

func ProtectedResourceIdFromContext(ctx context.Context) (ProtectedResourceId, bool) {
	id, ok := ctx.Value(protectedResourceIdContextKey{}).(ProtectedResourceId)
	return id, ok
}

func WithScopeId(ctx context.Context, id ScopeId) context.Context {
	return context.WithValue(ctx, scopeIdContextKey{}, id)
}

func ScopeIdFromContext(ctx context.Context) (ScopeId, bool) {
	id, ok := ctx.Value(scopeIdContextKey{}).(ScopeId)
	return id, ok
}
