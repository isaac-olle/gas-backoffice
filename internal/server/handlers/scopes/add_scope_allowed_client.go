package scopes

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	scopeAddClient "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/scopes/application/add_allowed_client"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/command"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

type AddScopeAllowedClientHandler struct {
	commandExecutor command.CommandExecutor
	scopeId         func(context.Context) (ids.ScopeId, bool)
	clientId        func(context.Context) (ids.ClientId, bool)
}

func NewAddScopeAllowedClientHandler(
	commandExecutor command.CommandExecutor,
	scopeId func(context.Context) (ids.ScopeId, bool),
	clientId func(context.Context) (ids.ClientId, bool),
) *AddScopeAllowedClientHandler {
	return &AddScopeAllowedClientHandler{commandExecutor: commandExecutor, scopeId: scopeId, clientId: clientId}
}

func (h *AddScopeAllowedClientHandler) Execute(c *gin.Context) (any, error) {
	scopeId, ok := h.scopeId(c.Request.Context())
	if !ok {
		return nil, app_errors.InternalError("MISSING_SCOPE_ID")
	}

	clientId, ok := h.clientId(c.Request.Context())
	if !ok {
		return nil, app_errors.InternalError("MISSING_CLIENT_ID")
	}

	return nil, h.commandExecutor.SingleExecution(
		c.Request.Context(), scopeAddClient.NewCommand(scopeId.String(), clientId.String()),
	)
}

func (h *AddScopeAllowedClientHandler) Success(c *gin.Context, _ any) error {
	c.Status(http.StatusNoContent)
	return nil
}

func (h *AddScopeAllowedClientHandler) Default(c *gin.Context) {
	c.Status(http.StatusInternalServerError)
}
