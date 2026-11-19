package protected_resources

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/command"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"

	prUpdateURL "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/protected_resources/application/update_url"
)

type UpdateProtectedResourceURLHandler struct {
	commandExecutor command.CommandExecutor
	prId            func(context.Context) (ids.ProtectedResourceId, bool)
}

func NewUpdateProtectedResourceURLHandler(
	commandExecutor command.CommandExecutor,
	prId func(context.Context) (ids.ProtectedResourceId, bool),
) *UpdateProtectedResourceURLHandler {
	return &UpdateProtectedResourceURLHandler{commandExecutor: commandExecutor, prId: prId}
}

type updateProtectedResourceURLRequest struct {
	URL *string `json:"url" binding:"required"`
}

func (h *UpdateProtectedResourceURLHandler) Execute(c *gin.Context) (any, error) {
	var req updateProtectedResourceURLRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, app_errors.InvalidRequestBodyError(err)
	}

	prId, ok := h.prId(c.Request.Context())
	if !ok {
		return nil, app_errors.InternalError("MISSING_PROTECTED_RESOURCE_ID")
	}

	return nil, h.commandExecutor.SingleExecution(c.Request.Context(), prUpdateURL.NewCommand(prId, *req.URL))
}

func (h *UpdateProtectedResourceURLHandler) Success(c *gin.Context, _ any) error {
	c.Status(http.StatusNoContent)
	return nil
}

func (h *UpdateProtectedResourceURLHandler) Default(c *gin.Context) {
	c.Status(http.StatusInternalServerError)
}
