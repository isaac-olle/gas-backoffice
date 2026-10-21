package protected_resources

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	prSoftDelete "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/protected_resources/application/soft_delete"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/command"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

type SoftDeleteProtectedResourceHandler struct {
	commandExecutor command.CommandExecutor
	prId            func(context.Context) (ids.ProtectedResourceId, bool)
}

func NewSoftDeleteProtectedResourceHandler(
	commandExecutor command.CommandExecutor,
	prId func(context.Context) (ids.ProtectedResourceId, bool),
) *SoftDeleteProtectedResourceHandler {
	return &SoftDeleteProtectedResourceHandler{commandExecutor: commandExecutor, prId: prId}
}

func (h *SoftDeleteProtectedResourceHandler) Execute(c *gin.Context) (any, error) {
	prId, ok := h.prId(c.Request.Context())
	if !ok {
		return nil, app_errors.InternalError("MISSING_PROTECTED_RESOURCE_ID")
	}

	return nil, h.commandExecutor.SingleExecution(c.Request.Context(), prSoftDelete.NewCommand(prId))
}

func (h *SoftDeleteProtectedResourceHandler) Success(c *gin.Context, _ any) error {
	c.Status(http.StatusNoContent)
	return nil
}

func (h *SoftDeleteProtectedResourceHandler) Default(c *gin.Context) {
	c.Status(http.StatusInternalServerError)
}
