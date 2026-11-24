package authenticators

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/auth"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

// AccessTokenCookieName is the HttpOnly cookie the backoffice login callback
// stores the access token in.
const AccessTokenCookieName = "access_token"
const RefreshTokenCookieName = "refresh_token"

// CookieAuthenticator reads the HttpOnly access cookie and verifies CSRF for writes.
type CookieAuthenticator struct {
	csrfVerifier auth.CSRFTokenVerifier
}

func NewCookieAuthenticator(csrfVerifier auth.CSRFTokenVerifier) *CookieAuthenticator {
	return &CookieAuthenticator{csrfVerifier: csrfVerifier}
}

func (a *CookieAuthenticator) VerifyCSRF(c *gin.Context, rawSessionID string) error {
	switch c.Request.Method {
	case http.MethodGet, http.MethodHead, http.MethodOptions:
		return nil
	}
	if a.csrfVerifier == nil {
		return nil
	}
	sessionID, err := ids.NewSessionIdFromString(rawSessionID)
	if err != nil {
		return app_errors.UnauthorizedError("INVALID_SESSION_ID")
	}
	return a.csrfVerifier.Verify(c.Request.Context(), sessionID, c.GetHeader("X-CSRF-Token"))
}

func (a *CookieAuthenticator) CanProcess(c *gin.Context) bool {
	_, err := c.Request.Cookie(AccessTokenCookieName)
	return err == nil
}

func (a *CookieAuthenticator) Token(c *gin.Context) (string, error) {
	token, err := c.Cookie(AccessTokenCookieName)
	if err != nil || token == "" {
		return "", app_errors.UnauthorizedError("MISSING_TOKEN").WithDescription("missing access token cookie")
	}
	return token, nil
}
