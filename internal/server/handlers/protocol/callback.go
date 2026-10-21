package protocol

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
)

type callbackRequest struct {
	Code  string `form:"code" binding:"required"`
	State string `form:"state" binding:"required"`
}

// CreateSessionExecutor is the narrow view of create_session.Service the HTTP
// layer needs, so tests can stub it without wiring the whole dependency graph.
type CreateSessionExecutor interface {
	CreateSession(ctx context.Context, code, state, idpId string) (string, error)
}

// CallbackHandler receives the external IDP redirect (?code&state) at
// /v1/oidc/callback/:idp_id, where :idp_id is the IDP reference name (e.g. google).
// It decodes the self-contained idp_state, exchanges the code with the IDP and
// creates the AuthorizationCodeSession already active, then redirects the
// browser back to the client's redirect_uri from /oidc/authorize.
type CallbackHandler struct {
	createSessionService CreateSessionExecutor
}

func NewCallbackHandler(createSessionService CreateSessionExecutor) *CallbackHandler {
	return &CallbackHandler{createSessionService: createSessionService}
}

func (h *CallbackHandler) Execute(c *gin.Context) (any, error) {
	var req callbackRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		return nil, app_errors.InvalidRequestBodyError(err)
	}

	return h.createSessionService.CreateSession(c.Request.Context(), req.Code, req.State, c.Param("idp_id"))
}

func (h *CallbackHandler) Success(c *gin.Context, result any) error {
	c.Redirect(http.StatusFound, result.(string))
	return nil
}

func (h *CallbackHandler) Default(c *gin.Context) {
	c.Status(http.StatusInternalServerError)
}
