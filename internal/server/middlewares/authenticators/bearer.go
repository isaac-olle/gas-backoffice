package authenticators

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
)

const bearerPrefix = "Bearer "

// BearerAuthenticator reads the access token from the Authorization header.
type BearerAuthenticator struct{}

func NewBearerAuthenticator() *BearerAuthenticator {
	return &BearerAuthenticator{}
}

func (a *BearerAuthenticator) Token(c *gin.Context) (string, error) {
	header := c.GetHeader("Authorization")
	token := strings.TrimPrefix(header, bearerPrefix)
	if !strings.HasPrefix(header, bearerPrefix) || token == "" {
		return "", app_errors.UnauthorizedError("MISSING_TOKEN").WithDescription("missing bearer token")
	}
	return token, nil
}
