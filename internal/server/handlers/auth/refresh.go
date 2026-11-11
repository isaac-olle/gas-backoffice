package auth

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	refreshSession "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/auth/application/refresh_session"
	"github.com/isaac-olle/gas-backoffice/internal/server/middlewares/authenticators"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/query"
)

type RefreshHandler struct{ queryBus query.Bus }

func NewRefreshHandler(queryBus query.Bus) *RefreshHandler {
	return &RefreshHandler{queryBus: queryBus}
}

func (h *RefreshHandler) Execute(c *gin.Context) (any, error) {
	refreshToken, err := c.Cookie(authenticators.RefreshTokenCookieName)
	if err != nil || refreshToken == "" {
		reason := "empty"
		var cause error = fmt.Errorf("cookie %q is empty", authenticators.RefreshTokenCookieName)
		if err != nil {
			reason = "missing"
			cause = fmt.Errorf("cannot read cookie %q: %w", authenticators.RefreshTokenCookieName, err)
		}
		return nil, app_errors.UnauthorizedError("MISSING_REFRESH_TOKEN").
			WithDescription("refresh_token cookie is missing or empty").
			WithCause(cause, true).
			WithDetails(map[string]any{"cookie_name": authenticators.RefreshTokenCookieName, "reason": reason})
	}
	return h.queryBus.Dispatch(c.Request.Context(), refreshSession.NewQuery(refreshToken))
}

func (h *RefreshHandler) Success(c *gin.Context, result any) error {
	res := result.(refreshSession.Result)
	c.SetCookie(authenticators.AccessTokenCookieName, res.AccessToken, 3600, "/", "", false, true)
	c.SetCookie(authenticators.RefreshTokenCookieName, res.RefreshToken, 30*24*60*60, "/", "", false, true)
	c.Status(http.StatusNoContent)
	return nil
}

func (h *RefreshHandler) Default(c *gin.Context) { c.Status(http.StatusInternalServerError) }
