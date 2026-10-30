package auth

import (
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
		return nil, app_errors.UnauthorizedError("MISSING_REFRESH_TOKEN")
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
