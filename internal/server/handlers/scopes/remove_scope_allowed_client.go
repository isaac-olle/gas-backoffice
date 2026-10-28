package scopes

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	scopeRemoveClient "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/scopes/application/remove_allowed_client"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/command"
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

type RemoveScopeAllowedClientHandler struct {
	commandExecutor command.CommandExecutor
	scopeId         func(context.Context) (ids2.ScopeId, bool)
	clientId        func(context.Context) (ids2.ClientId, bool)
}

func NewRemoveScopeAllowedClientHandler(
	commandExecutor command.CommandExecutor,
	scopeId func(context.Context) (ids2.ScopeId, bool),
	clientId func(context.Context) (ids2.ClientId, bool),
) *RemoveScopeAllowedClientHandler {
	return &RemoveScopeAllowedClientHandler{commandExecutor: commandExecutor, scopeId: scopeId, clientId: clientId}
}

func (h *RemoveScopeAllowedClientHandler) Execute(c *gin.Context) (any, error) {
	scopeId, ok := h.scopeId(c.Request.Context())
	if !ok {
		return nil, app_errors.InternalError("MISSING_SCOPE_ID")
	}

	clientId, ok := h.clientId(c.Request.Context())
	if !ok {
		return nil, app_errors.InternalError("MISSING_CLIENT_ID")
	}

	return nil, h.commandExecutor.SingleExecution(
		c.Request.Context(), scopeRemoveClient.NewCommand(scopeId.String(), clientId.String()),
	)
}

func (h *RemoveScopeAllowedClientHandler) Success(c *gin.Context, _ any) error {
	c.Status(http.StatusNoContent)
	return nil
}

func (h *RemoveScopeAllowedClientHandler) Default(c *gin.Context) {
	c.Status(http.StatusInternalServerError)
}
