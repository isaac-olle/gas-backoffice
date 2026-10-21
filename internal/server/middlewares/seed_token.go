package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/isaac-olle/gas-backoffice/internal/server/middlewares/authenticators"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
)

func SeedToken(seedToken string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if seedToken == "" {
			_ = c.Error(app_errors.UnauthorizedError("SEED_TOKEN_NOT_CONFIGURED"))
			c.Abort()
			return
		}
		token, err := authenticators.NewBearerAuthenticator().Token(c)
		if err != nil || token != seedToken {
			_ = c.Error(app_errors.UnauthorizedError("INVALID_SEED_TOKEN"))
			c.Abort()
			return
		}
		c.Next()
	}
}
