package protected_resources

import (
	"context"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/command"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"net/http"

	"github.com/gin-gonic/gin"
	prUpdateName "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/protected_resources/application/update_name"
)

type UpdateProtectedResourceNameHandler struct {
	commandBus command.Bus
	prId       func(context.Context) (ids.ProtectedResourceId, bool)
}

func NewUpdateProtectedResourceNameHandler(
	commandBus command.Bus,
	prId func(context.Context) (ids.ProtectedResourceId, bool),
) *UpdateProtectedResourceNameHandler {
	return &UpdateProtectedResourceNameHandler{commandBus: commandBus, prId: prId}
}

type updateProtectedResourceNameRequest struct {
	Name string `json:"name" binding:"required"`
}

func (h *UpdateProtectedResourceNameHandler) Execute(c *gin.Context) (any, error) {
	var req updateProtectedResourceNameRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, app_errors.InvalidRequestBodyError(err)
	}

	prId, ok := h.prId(c.Request.Context())
	if !ok {
		return nil, app_errors.InternalError("MISSING_PROTECTED_RESOURCE_ID")
	}

	return nil, h.commandBus.Dispatch(c.Request.Context(), prUpdateName.NewCommand(prId, req.Name))
}

func (h *UpdateProtectedResourceNameHandler) Success(c *gin.Context, _ any) error {
	c.Status(http.StatusNoContent)
	return nil
}

func (h *UpdateProtectedResourceNameHandler) Default(c *gin.Context) {
	c.Status(http.StatusInternalServerError)
}
