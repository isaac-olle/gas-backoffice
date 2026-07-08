package users

import (
	"context"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/command"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"net/http"

	"github.com/gin-gonic/gin"
	userChangeName "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/users/application/change_name"
)

type ChangeUserNameHandler struct {
	commandBus command.Bus
	userId     func(context.Context) (ids.UserId, bool)
}

func NewChangeUserNameHandler(
	commandBus command.Bus,
	userId func(context.Context) (ids.UserId, bool),
) *ChangeUserNameHandler {
	return &ChangeUserNameHandler{commandBus: commandBus, userId: userId}
}

type changeUserNameRequest struct {
	Name string `json:"name" binding:"required"`
}

func (h *ChangeUserNameHandler) Execute(c *gin.Context) (any, error) {
	userId, ok := h.userId(c.Request.Context())
	if !ok {
		return nil, app_errors.InternalError("MISSING_USER_ID")
	}

	var req changeUserNameRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, err
	}

	return nil, h.commandBus.Dispatch(c.Request.Context(), userChangeName.NewCommand(userId, req.Name))
}

func (h *ChangeUserNameHandler) Success(c *gin.Context, _ any) error {
	c.Status(http.StatusNoContent)
	return nil
}

func (h *ChangeUserNameHandler) Default(c *gin.Context) {
	c.Status(http.StatusInternalServerError)
}
