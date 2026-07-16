package scopes

import (
	"context"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/command"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"net/http"

	"github.com/gin-gonic/gin"
	scopeAddClient "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/scopes/application/add_allowed_client"
)

type AddScopeAllowedClientHandler struct {
	commandBus command.Bus
	scopeId    func(context.Context) (ids.ScopeId, bool)
}

func NewAddScopeAllowedClientHandler(
	commandBus command.Bus,
	scopeId func(context.Context) (ids.ScopeId, bool),
) *AddScopeAllowedClientHandler {
	return &AddScopeAllowedClientHandler{commandBus: commandBus, scopeId: scopeId}
}

type addScopeAllowedClientRequest struct {
	ClientID string `json:"client_id" binding:"required"`
}

func (h *AddScopeAllowedClientHandler) Execute(c *gin.Context) (any, error) {
	var req addScopeAllowedClientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, app_errors.InvalidRequestBodyError(err)
	}

	scopeId, ok := h.scopeId(c.Request.Context())
	if !ok {
		return nil, app_errors.InternalError("MISSING_SCOPE_ID")
	}

	return nil, h.commandBus.Dispatch(c.Request.Context(), scopeAddClient.NewCommand(scopeId.String(), req.ClientID))
}

func (h *AddScopeAllowedClientHandler) Success(c *gin.Context, _ any) error {
	c.Status(http.StatusNoContent)
	return nil
}

func (h *AddScopeAllowedClientHandler) Default(c *gin.Context) {
	c.Status(http.StatusInternalServerError)
}
