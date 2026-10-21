package protocol

import (
	"context"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session/application/token/flows"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session/application/token"
)

// TokenExecutor is the narrow view of token.Service the HTTP layer needs, so
// tests can stub it without wiring the whole dependency graph.
type TokenExecutor interface {
	Execute(ctx context.Context, grantType string, req token.Request) (*session.TokenResponse, error)
}

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
	tokenService TokenExecutor
}

func NewTokenHandler(tokenService TokenExecutor) *TokenHandler {
	return &TokenHandler{tokenService: tokenService}
}

func (h *TokenHandler) Execute(c *gin.Context) (any, error) {
	clientID, _ := AuthenticatedClientIDFromContext(c.Request.Context())
	isPublic := IsPublicClientFromContext(c.Request.Context())
	ctx := c.Request.Context()

	var grantReq grantTypeRequest
	if err := c.ShouldBind(&grantReq); err != nil {
		return nil, app_errors.AddTrace(err)
	}

	req := token.Request{ClientID: clientID, IsPublic: isPublic}

	switch grantReq.GrantType {
	case flows.AuthorizationCodeGrantType:
		var authReq authorizationCodeRequest
		if err := c.ShouldBind(&authReq); err != nil {
			return nil, app_errors.AddTrace(err)
		}
		req.Code = authReq.Code
		req.RedirectURI = authReq.RedirectURI
		req.CodeVerifier = authReq.CodeVerifier

	case flows.ClientCredentialsGrantType:
		var ccReq clientCredentialsRequest
		if err := c.ShouldBind(&ccReq); err != nil {
			return nil, app_errors.AddTrace(err)
		}
		req.Scope = ccReq.Scope
	}

	return h.tokenService.Execute(ctx, grantReq.GrantType, req)
}

func (h *TokenHandler) Success(c *gin.Context, result any) error {
	resp := result.(*session.TokenResponse)
	response := gin.H{
		"access_token": resp.AccessToken,
		"token_type":   resp.TokenType,
		"expires_in":   resp.ExpiresIn,
		"scope":        resp.Scope,
	}
	if resp.IDToken != "" {
		response["id_token"] = resp.IDToken
	}
	c.JSON(http.StatusOK, response)
	return nil
}

func (h *TokenHandler) Default(c *gin.Context) {
	c.Status(http.StatusInternalServerError)
}
