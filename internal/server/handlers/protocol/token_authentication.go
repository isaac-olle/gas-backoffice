package protocol

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/isaac-olle/gas-backoffice/internal/server/handlers/protocol/authentication_methods"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
)

type tokenAuthContextKey struct{}

func AuthenticatedClientIDFromContext(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(tokenAuthContextKey{}).(string)
	return id, ok
}

type TokenAuthenticationHandler struct {
	factory *authentication_methods.ClientAuthMethodFactory
}

func NewTokenAuthenticationHandler(factory *authentication_methods.ClientAuthMethodFactory) *TokenAuthenticationHandler {
	return &TokenAuthenticationHandler{factory: factory}
}

func (h *TokenAuthenticationHandler) Execute(c *gin.Context) (any, error) {
	clientID, err := h.factory.Process(c)
	if err != nil {
		return nil, err
	}
	ctx := context.WithValue(c.Request.Context(), tokenAuthContextKey{}, clientID)
	c.Request = c.Request.WithContext(ctx)
	return nil, nil
}

func (h *TokenAuthenticationHandler) Success(c *gin.Context, _ any) error {
	c.Next()
	return nil
}

func (h *TokenAuthenticationHandler) Default(c *gin.Context) {
	_ = c.Error(app_errors.UnauthorizedError("UNAUTHORIZED"))
	c.Abort()
}
