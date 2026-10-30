package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	logoutSession "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/auth/application/logout"
	verifyJwt "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/key_pairs/application/verify_jwt"
	"github.com/isaac-olle/gas-backoffice/internal/server/middlewares/authenticators"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/command"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/query"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

type LogoutHandler struct {
	commands command.CommandExecutor
	queries  query.Bus
	issuer   string
}

func NewLogoutHandler(commands command.CommandExecutor, queries query.Bus, issuer string) *LogoutHandler {
	return &LogoutHandler{commands: commands, queries: queries, issuer: issuer}
}

func (h *LogoutHandler) Execute(c *gin.Context) (any, error) {
	if sessionID, ok := h.sessionID(c, authenticators.AccessTokenCookieName, ""); ok {
		return nil, h.commands.SingleExecution(c.Request.Context(), logoutSession.NewCommand(sessionID))
	} else if sessionID, ok := h.sessionID(c, authenticators.RefreshTokenCookieName, "backoffice_refresh"); ok {
		return nil, h.commands.SingleExecution(c.Request.Context(), logoutSession.NewCommand(sessionID))
	}
	return nil, nil
}

func (h *LogoutHandler) sessionID(c *gin.Context, cookieName, tokenType string) (ids.SessionId, bool) {
	token, err := c.Cookie(cookieName)
	if err != nil || token == "" {
		return ids.SessionId{}, false
	}
	qry := verifyJwt.NewQuery(token).
		WithExpectedClaim("iss", h.issuer).
		WithExpectedClaim("aud", h.issuer).
		WithClaimsToReturn(map[string]any{"sid": nil})
	if tokenType != "" {
		qry.WithExpectedClaim("typ", tokenType)
	}
	claims, err := query.DispatchTyped[map[string]any](c.Request.Context(), h.queries, qry)
	if err != nil {
		return ids.SessionId{}, false
	}
	raw, ok := claims["sid"].(string)
	if !ok {
		return ids.SessionId{}, false
	}
	sessionID, err := ids.NewSessionIdFromString(raw)
	return sessionID, err == nil
}

func (h *LogoutHandler) Success(c *gin.Context, _ any) error {
	c.SetCookie(authenticators.AccessTokenCookieName, "", -1, "/", "", false, true)
	c.SetCookie(authenticators.RefreshTokenCookieName, "", -1, "/", "", false, true)
	c.Status(http.StatusNoContent)
	return nil
}

func (h *LogoutHandler) Default(c *gin.Context) { c.Status(http.StatusNoContent) }
