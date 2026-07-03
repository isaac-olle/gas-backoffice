package scopes

import (
	"context"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/command"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"net/http"

	"github.com/gin-gonic/gin"
	scopeSoftDelete "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/scopes/application/soft_delete"
)

type SoftDeleteScopeHandler struct {
	commandBus command.Bus
	scopeId    func(context.Context) (ids.ScopeId, bool)
}

func NewSoftDeleteScopeHandler(
	commandBus command.Bus,
	scopeId func(context.Context) (ids.ScopeId, bool),
) *SoftDeleteScopeHandler {
	return &SoftDeleteScopeHandler{commandBus: commandBus, scopeId: scopeId}
}

func (h *SoftDeleteScopeHandler) Execute(c *gin.Context) (any, error) {
	scopeId, ok := h.scopeId(c.Request.Context())
	if !ok {
		return nil, app_errors.InternalError("MISSING_SCOPE_ID")
	}
	return nil, h.commandBus.Dispatch(c.Request.Context(), scopeSoftDelete.NewCommand(scopeId.String()))
}

func (h *SoftDeleteScopeHandler) Success(c *gin.Context, _ any) error {
	c.Status(http.StatusNoContent)
	return nil
}

func (h *SoftDeleteScopeHandler) Default(c *gin.Context) {
	c.Status(http.StatusInternalServerError)
}
