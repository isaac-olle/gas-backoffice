package protocol

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/protocol/session/application/authorize"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
)

// AuthorizeExecutor is the narrow view of authorize.Service the HTTP layer
// needs, so tests can stub it without wiring the whole dependency graph.
type AuthorizeExecutor interface {
	Authorize(ctx context.Context, req authorize.Request) (string, error)
}

type authorizeRequest struct {
	ResponseType        string `form:"response_type" binding:"required"`
	ClientID            string `form:"client_id" binding:"required"`
	RedirectURI         string `form:"redirect_uri" binding:"required"`
	Scope               string `form:"scope" binding:"required"`
	State               string `form:"state"`
	CodeChallenge       string `form:"code_challenge"`
	CodeChallengeMethod string `form:"code_challenge_method"`
	IdpReferenceName    string `form:"idp_reference_name"`
}

type AuthorizeHandler struct {
	authorizeService AuthorizeExecutor
}

func NewAuthorizeHandler(authorizeService AuthorizeExecutor) *AuthorizeHandler {
	return &AuthorizeHandler{authorizeService: authorizeService}
}

func (h *AuthorizeHandler) Execute(c *gin.Context) (any, error) {
	var req authorizeRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		return nil, app_errors.InvalidRequestBodyError(err)
	}

	return h.authorizeService.Authorize(c.Request.Context(), authorize.Request{
		ResponseType:        req.ResponseType,
		ClientID:            req.ClientID,
		RedirectURI:         req.RedirectURI,
		Scope:               req.Scope,
		State:               req.State,
		CodeChallenge:       req.CodeChallenge,
		CodeChallengeMethod: req.CodeChallengeMethod,
		IdpReferenceName:    req.IdpReferenceName,
	})
}

func (h *AuthorizeHandler) Success(c *gin.Context, result any) error {
	c.Redirect(http.StatusFound, result.(string))
	return nil
}

func (h *AuthorizeHandler) Default(c *gin.Context) {
	c.Status(http.StatusInternalServerError)
}
