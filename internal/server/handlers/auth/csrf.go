package auth

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/auth"
	verifyJWT "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/key_pairs/application/verify_jwt"
	"github.com/isaac-olle/gas-backoffice/internal/server/middlewares/authenticators"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/query"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

type CSRFTokenGetter interface {
	Get(context.Context, ids.SessionId) (string, error)
}

type CSRFHandler struct {
	tokens          CSRFTokenGetter
	refreshVerifier auth.RefreshTokenVerifier
	queries         query.Bus
	issuer          string
}

func NewCSRFHandler(tokens CSRFTokenGetter, refreshVerifier auth.RefreshTokenVerifier, queries query.Bus, issuer string) *CSRFHandler {
	return &CSRFHandler{tokens: tokens, refreshVerifier: refreshVerifier, queries: queries, issuer: issuer}
}

func (h *CSRFHandler) Execute(c *gin.Context) (any, error) {
	ctx := c.Request.Context()
	if id, ok := h.accessSessionID(c); ok {
		return h.tokens.Get(ctx, id)
	}
	token, err := c.Cookie(authenticators.RefreshTokenCookieName)
	if err == nil && token != "" {
		claims, err := h.refreshVerifier.Verify(ctx, token)
		if err == nil {
			return h.tokens.Get(ctx, claims.SessionID)
		}
	}
	return nil, app_errors.UnauthorizedError("INVALID_SESSION_TOKEN")
}

func (h *CSRFHandler) accessSessionID(c *gin.Context) (ids.SessionId, bool) {
	token, err := c.Cookie(authenticators.AccessTokenCookieName)
	if err != nil || token == "" {
		return ids.SessionId{}, false
	}
	claims, err := query.DispatchTyped[map[string]any](c.Request.Context(), h.queries,
		verifyJWT.NewQuery(token).WithExpectedClaim("iss", h.issuer).
			WithExpectedClaim("aud", h.issuer).WithVerifyExpiration(),
	)
	if err != nil || claims["typ"] == "backoffice_refresh" {
		return ids.SessionId{}, false
	}
	rawID, _ := claims["sid"].(string)
	id, err := ids.NewSessionIdFromString(rawID)
	return id, err == nil
}

func (h *CSRFHandler) Success(c *gin.Context, result any) error {
	c.Header("Cache-Control", "no-store")
	c.JSON(http.StatusOK, gin.H{"csrf_token": result})
	return nil
}

func (h *CSRFHandler) Default(c *gin.Context) {
	c.Status(http.StatusInternalServerError)
}
