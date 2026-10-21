package persistance

import (
	"context"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

// RequestCachedSessionRepository wraps a session.Repository and, when the
// context carries a cache (see WithCache), remembers the first FindById/
// FindByCode lookup so a second lookup for the same key within the same
// request (e.g. a command followed by a query, both triggered by one HTTP
// handler) is served from memory instead of hitting the database again.
// Since sessions are mutated in place, a cached pointer already reflects any
// mutation made after the first lookup — no invalidation is needed. Falls
// through to the inner repository untouched when no cache is present.
type RequestCachedSessionRepository struct {
	inner session.Repository
}

func NewRequestCachedSessionRepository(inner session.Repository) *RequestCachedSessionRepository {
	return &RequestCachedSessionRepository{inner: inner}
}

type sessionCacheKey struct{}

type sessionCache struct {
	byId   map[string]*session.AuthorizationCodeSession
	byCode map[string]*session.AuthorizationCodeSession
}

// WithCache returns a new context carrying an empty cache. Pass the resulting
// context to every command/query dispatched within the same request so their
// repository lookups can share it.
func WithCache(ctx context.Context) context.Context {
	return context.WithValue(ctx, sessionCacheKey{}, &sessionCache{
		byId:   make(map[string]*session.AuthorizationCodeSession),
		byCode: make(map[string]*session.AuthorizationCodeSession),
	})
}

func (r *RequestCachedSessionRepository) Save(ctx context.Context, sess *session.AuthorizationCodeSession) error {
	return r.inner.Save(ctx, sess)
}

func (r *RequestCachedSessionRepository) FindById(ctx context.Context, id ids.SessionId) (*session.AuthorizationCodeSession, error) {
	cache, ok := ctx.Value(sessionCacheKey{}).(*sessionCache)
	if !ok {
		return r.inner.FindById(ctx, id)
	}
	if sess, found := cache.byId[id.String()]; found {
		return sess, nil
	}
	sess, err := r.inner.FindById(ctx, id)
	if err == nil && sess != nil {
		cache.byId[id.String()] = sess
	}
	return sess, err
}

func (r *RequestCachedSessionRepository) FindByCode(ctx context.Context, code string) (*session.AuthorizationCodeSession, error) {
	cache, ok := ctx.Value(sessionCacheKey{}).(*sessionCache)
	if !ok {
		return r.inner.FindByCode(ctx, code)
	}
	if sess, found := cache.byCode[code]; found {
		return sess, nil
	}
	sess, err := r.inner.FindByCode(ctx, code)
	if err == nil && sess != nil {
		cache.byCode[code] = sess
	}
	return sess, err
}

func (r *RequestCachedSessionRepository) SoftDelete(ctx context.Context, id ids.SessionId) error {
	return r.inner.SoftDelete(ctx, id)
}

func (r *RequestCachedSessionRepository) HardDelete(ctx context.Context, id ids.SessionId) error {
	return r.inner.HardDelete(ctx, id)
}
