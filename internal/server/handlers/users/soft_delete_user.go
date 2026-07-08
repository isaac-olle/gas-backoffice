package users

import (
	"context"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/command"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"net/http"

	"github.com/gin-gonic/gin"
	userSoftDelete "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/users/application/soft_delete"
)

type SoftDeleteUserHandler struct {
	commandBus command.Bus
	userId     func(context.Context) (ids.UserId, bool)
}

func NewSoftDeleteUserHandler(
	commandBus command.Bus,
	userId func(context.Context) (ids.UserId, bool),
) *SoftDeleteUserHandler {
	return &SoftDeleteUserHandler{commandBus: commandBus, userId: userId}
}

func (h *SoftDeleteUserHandler) Execute(c *gin.Context) (any, error) {
	userId, ok := h.userId(c.Request.Context())
	if !ok {
		return nil, app_errors.InternalError("MISSING_USER_ID")
	}
	return nil, h.commandBus.Dispatch(c.Request.Context(), userSoftDelete.NewCommand(userId))
}

func (h *SoftDeleteUserHandler) Success(c *gin.Context, _ any) error {
	c.Status(http.StatusNoContent)
	return nil
}

func (h *SoftDeleteUserHandler) Default(c *gin.Context) {
	c.Status(http.StatusInternalServerError)
}
