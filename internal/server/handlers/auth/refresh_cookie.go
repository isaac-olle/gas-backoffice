package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/isaac-olle/gas-backoffice/internal/server/middlewares/authenticators"
)

const refreshCookiePath = "/v1/backoffice/auth"

func setRefreshCookie(c *gin.Context, token string, maxAge int) {
	c.SetCookie(authenticators.RefreshTokenCookieName, token, maxAge, refreshCookiePath, "", true, true)
	// Clear the cookie issued by deployments that used the root path.
	c.SetCookie(authenticators.RefreshTokenCookieName, "", -1, "/", "", true, true)
}
