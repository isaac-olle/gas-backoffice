package scopes

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	scopeDisable "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/scopes/application/disable"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/command"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

type DisableScopeHandler struct {
	commandExecutor command.CommandExecutor
	scopeId         func(context.Context) (ids.ScopeId, bool)
}

func NewDisableScopeHandler(
	commandExecutor command.CommandExecutor,
	scopeId func(context.Context) (ids.ScopeId, bool),
) *DisableScopeHandler {
	return &DisableScopeHandler{commandExecutor: commandExecutor, scopeId: scopeId}
}

func (h *DisableScopeHandler) Execute(c *gin.Context) (any, error) {
	scopeId, ok := h.scopeId(c.Request.Context())
	if !ok {
		return nil, app_errors.InternalError("MISSING_SCOPE_ID")
	}
	return nil, h.commandExecutor.SingleExecution(
		c.Request.Context(), scopeDisable.NewCommand(scopeId.String()),
	)
}

func (h *DisableScopeHandler) Success(c *gin.Context, _ any) error {
	c.Status(http.StatusNoContent)
	return nil
}

func (h *DisableScopeHandler) Default(c *gin.Context) {
	c.Status(http.StatusInternalServerError)
}
