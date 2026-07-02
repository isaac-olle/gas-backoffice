package protocol

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/isaac-olle/gas/internal/bounded_contexts/protocol/session"
	exchangeCode "github.com/isaac-olle/gas/internal/bounded_contexts/protocol/session/application/exchange_code"
	issueCC "github.com/isaac-olle/gas/internal/bounded_contexts/protocol/session/application/issue_client_credentials"
	"github.com/isaac-olle/gas/internal/shared/buses/query"
)

type grantTypeRequest struct {
	GrantType string `form:"grant_type" binding:"required"`
}

type authorizationCodeRequest struct {
	Code         string `form:"code" binding:"required"`
	RedirectURI  string `form:"redirect_uri" binding:"required"`
	CodeVerifier string `form:"code_verifier"`
}

type clientCredentialsRequest struct {
	Scope string `form:"scope" binding:"required"`
}

type TokenHandler struct {
	queryBus query.Bus
}

func NewTokenHandler(qryBus query.Bus) *TokenHandler {
	return &TokenHandler{queryBus: qryBus}
}

func (h *TokenHandler) Execute(c *gin.Context) (any, error) {
	clientID, _ := AuthenticatedClientIDFromContext(c.Request.Context())

	var grantReq grantTypeRequest
	if err := c.ShouldBind(&grantReq); err != nil {
		return nil, err
	}

	switch grantReq.GrantType {
	case "authorization_code":
		var req authorizationCodeRequest
		if err := c.ShouldBind(&req); err != nil {
			return nil, err
		}
		return query.DispatchTyped[*session.TokenResponse](
			c.Request.Context(), h.queryBus,
			exchangeCode.NewQuery(req.Code, clientID, req.RedirectURI, req.CodeVerifier),
		)

	case "client_credentials":
		var req clientCredentialsRequest
		if err := c.ShouldBind(&req); err != nil {
			return nil, err
		}
		return query.DispatchTyped[*session.TokenResponse](
			c.Request.Context(), h.queryBus,
			issueCC.NewQuery(clientID, req.Scope),
		)

	default:
		return nil, unsupportedGrantTypeError(grantReq.GrantType)
	}
}

func (h *TokenHandler) Success(c *gin.Context, result any) error {
	resp := result.(*session.TokenResponse)
	c.JSON(http.StatusOK, gin.H{
		"access_token": resp.Token,
		"token_type":   "Bearer",
		"expires_in":   3600,
		"scope":        resp.Scope,
	})
	return nil
}

func (h *TokenHandler) Default(c *gin.Context) {
	c.Status(http.StatusInternalServerError)
}
