package ids

import "context"

type adminIdContextKey struct{}
type tenantIdContextKey struct{}
type clientIdContextKey struct{}
type protectedResourceIdContextKey struct{}
type scopeIdContextKey struct{}

func WithAdminId(ctx context.Context, id AdminId) context.Context {
	return context.WithValue(ctx, adminIdContextKey{}, id)
}

func AdminIdFromContext(ctx context.Context) (AdminId, bool) {
	id, ok := ctx.Value(adminIdContextKey{}).(AdminId)
	return id, ok
}

func WithTenantId(ctx context.Context, id TenantId) context.Context {
	return context.WithValue(ctx, tenantIdContextKey{}, id)
}

func TenantIdFromContext(ctx context.Context) (TenantId, bool) {
	id, ok := ctx.Value(tenantIdContextKey{}).(TenantId)
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
