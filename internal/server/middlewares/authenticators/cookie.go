package authenticators

import (
	"github.com/gin-gonic/gin"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
)

// AccessTokenCookieName is the HttpOnly cookie the backoffice login callback
// stores the access token in.
const AccessTokenCookieName = "access_token"

// CookieAuthenticator reads the access token from the HttpOnly cookie the
// browser sends along with requests coming from the backoffice frontend.
type CookieAuthenticator struct{}

func NewCookieAuthenticator() *CookieAuthenticator {
	return &CookieAuthenticator{}
}

func (a *CookieAuthenticator) Token(c *gin.Context) (string, error) {
	token, err := c.Cookie(AccessTokenCookieName)
	if err != nil || token == "" {
		return "", app_errors.UnauthorizedError("MISSING_TOKEN").WithDescription("missing access token cookie")
	}
	return token, nil
}
