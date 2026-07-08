package users

import (
	"context"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/command"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"net/http"

	"github.com/gin-gonic/gin"
	userHardDelete "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/users/application/hard_delete"
)

type HardDeleteUserHandler struct {
	commandBus command.Bus
	userId     func(context.Context) (ids.UserId, bool)
}

func NewHardDeleteUserHandler(
	commandBus command.Bus,
	userId func(context.Context) (ids.UserId, bool),
) *HardDeleteUserHandler {
	return &HardDeleteUserHandler{commandBus: commandBus, userId: userId}
}

func (h *HardDeleteUserHandler) Execute(c *gin.Context) (any, error) {
	userId, ok := h.userId(c.Request.Context())
	if !ok {
		return nil, app_errors.InternalError("MISSING_USER_ID")
	}
	return nil, h.commandBus.Dispatch(c.Request.Context(), userHardDelete.NewCommand(userId.String()))
}

func (h *HardDeleteUserHandler) Success(c *gin.Context, _ any) error {
	c.Status(http.StatusNoContent)
	return nil
}

func (h *HardDeleteUserHandler) Default(c *gin.Context) {
	c.Status(http.StatusInternalServerError)
}
