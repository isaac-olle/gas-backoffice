package protocol

import (
	"context"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session/application/authenticate_client"
)

type tokenClientIDKey struct{}
type tokenIsPublicKey struct{}

func AuthenticatedClientIDFromContext(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(tokenClientIDKey{}).(string)
	return id, ok
}

func IsPublicClientFromContext(ctx context.Context) bool {
	isPublic, _ := ctx.Value(tokenIsPublicKey{}).(bool)
	return isPublic
}

func WithAuthenticatedClient(ctx context.Context, clientID string, isPublic bool) context.Context {
	ctx = context.WithValue(ctx, tokenClientIDKey{}, clientID)
	ctx = context.WithValue(ctx, tokenIsPublicKey{}, isPublic)
	return ctx
}

// AuthenticateClientExecutor is the narrow view of authenticate_client.Service
// the HTTP layer needs, so tests can stub it without wiring the whole
// dependency graph.
type AuthenticateClientExecutor interface {
	Authenticate(ctx context.Context, creds authenticate_client.Credentials) (authenticate_client.Result, error)
}

func TokenAuthentication(authService AuthenticateClientExecutor) gin.HandlerFunc {
	return func(c *gin.Context) {
		creds := authenticate_client.Credentials{
			AuthorizationHeader: c.GetHeader("Authorization"),
			IsFormEncoded:       strings.Contains(c.GetHeader("Content-Type"), "application/x-www-form-urlencoded"),
			ClientID:            c.PostForm("client_id"),
			ClientSecret:        c.PostForm("client_secret"),
			ClientAssertionType: c.PostForm("client_assertion_type"),
			ClientAssertion:     c.PostForm("client_assertion"),
		}

		result, err := authService.Authenticate(c.Request.Context(), creds)
		if err != nil {
			_ = c.Error(err)
			c.Abort()
			return
		}
		c.Request = c.Request.WithContext(
			WithAuthenticatedClient(c.Request.Context(), result.ClientID, result.IsPublic),
		)
		c.Next()
	}
}
