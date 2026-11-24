package authenticators

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
)

// BearerAuthenticator reads the access token from the Authorization header.
type BearerAuthenticator struct{}

func NewBearerAuthenticator() *BearerAuthenticator {
	return &BearerAuthenticator{}
}

func (a *BearerAuthenticator) CanProcess(c *gin.Context) bool {
	return len(c.Request.Header.Values("Authorization")) > 0
}

func (a *BearerAuthenticator) Token(c *gin.Context) (string, error) {
	headers := c.Request.Header.Values("Authorization")
	if len(headers) == 0 {
		return "", app_errors.UnauthorizedError("MISSING_TOKEN").WithDescription("missing bearer token")
	}
	parts := strings.Fields(headers[0])
	if len(headers) != 1 || len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return "", app_errors.UnauthorizedError("INVALID_AUTHORIZATION_HEADER").WithDescription("expected Authorization: Bearer <token>")
	}
	return parts[1], nil
}
