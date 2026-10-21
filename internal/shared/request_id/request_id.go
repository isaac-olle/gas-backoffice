package request_id

import "context"

// Header is the HTTP header used to carry the request id across process
// boundaries, both inbound (from a client) and outbound (to an upstream
// endpoint we call on the same request's behalf).
const Header = "X-Request-ID"

type contextKey struct{}

// WithContext attaches the request id to ctx so it survives past the
// gin.Context it was read from (e.g. into an outbound HTTP call made from an
// application service).
func WithContext(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, contextKey{}, requestID)
}

// FromContext returns the request id stored by WithContext, or "" if none.
func FromContext(ctx context.Context) string {
	requestID, _ := ctx.Value(contextKey{}).(string)
	return requestID
}
