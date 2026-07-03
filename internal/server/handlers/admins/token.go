package admins

import (
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/query"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/actors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	verifyAssertion "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/admin_assertion_keys/application/verify_assertion"
	generateJwt "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/key_pairs/application/generate_jwt"
)

type TokenHandler struct {
	queryBus     query.Bus
	signingKeyId string
	jwtIssuer    string
}

func NewTokenHandler(queryBus query.Bus, signingKeyId, jwtIssuer string) *TokenHandler {
	return &TokenHandler{queryBus: queryBus, signingKeyId: signingKeyId, jwtIssuer: jwtIssuer}
}

type adminTokenRequest struct {
	GrantType string `form:"grant_type" binding:"required"`
	Assertion string `form:"assertion" binding:"required"`
}

func (h *TokenHandler) Execute(c *gin.Context) (any, error) {
	var req adminTokenRequest
	if err := c.ShouldBind(&req); err != nil {
		return nil, err
	}

	adminActor, err := query.DispatchTyped[*actors.AdminActor](
		c.Request.Context(),
		h.queryBus,
		verifyAssertion.NewQuery(req.Assertion, req.GrantType),
	)
	if err != nil {
		return nil, err
	}

	p := adminActor.ToPrimitives()

	tokenString, err := query.DispatchTyped[string](
		c.Request.Context(),
		h.queryBus,
		generateJwt.NewQuery(
			h.signingKeyId, map[string]any{
				"iss":  h.jwtIssuer,
				"sub":  p.Id,
				"aud":  h.jwtIssuer,
				"role": p.Role,
			}, 24*time.Hour),
	)
	if err != nil {
		return nil, err
	}

	return gin.H{
		"access_token": tokenString,
		"token_type":   "Bearer",
		"expires_in":   86400,
	}, nil
}

func (h *TokenHandler) Success(c *gin.Context, result any) error {
	c.JSON(http.StatusOK, result)
	return nil
}

func (h *TokenHandler) Default(c *gin.Context) {
	c.Status(http.StatusInternalServerError)
}
