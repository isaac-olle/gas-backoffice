package scopes

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	scopeEnable "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/scopes/application/enable"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/command"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

type EnableScopeHandler struct {
	commandExecutor command.CommandExecutor
	scopeId         func(context.Context) (ids.ScopeId, bool)
}

func NewEnableScopeHandler(
	commandExecutor command.CommandExecutor,
	scopeId func(context.Context) (ids.ScopeId, bool),
) *EnableScopeHandler {
	return &EnableScopeHandler{commandExecutor: commandExecutor, scopeId: scopeId}
}

func (h *EnableScopeHandler) Execute(c *gin.Context) (any, error) {
	scopeId, ok := h.scopeId(c.Request.Context())
	if !ok {
		return nil, app_errors.InternalError("MISSING_SCOPE_ID")
	}
	return nil, h.commandExecutor.SingleExecution(
		c.Request.Context(), scopeEnable.NewCommand(scopeId.String()),
	)
}

func (h *EnableScopeHandler) Success(c *gin.Context, _ any) error {
	c.Status(http.StatusNoContent)
	return nil
}

func (h *EnableScopeHandler) Default(c *gin.Context) {
	c.Status(http.StatusInternalServerError)
}
