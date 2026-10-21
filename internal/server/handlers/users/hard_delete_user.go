package users

import (
	"context"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/command"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/authors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"net/http"

	"github.com/gin-gonic/gin"
	userHardDelete "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/tenants/application/hard_delete"
)

type HardDeleteUserHandler struct {
	commandExecutor command.CommandExecutor
	userId          func(context.Context) (ids.UserId, bool)
	author          func(context.Context) (*authors.UserAuthor, bool)
}

func NewHardDeleteUserHandler(
	commandExecutor command.CommandExecutor,
	userId func(context.Context) (ids.UserId, bool),
	author func(context.Context) (*authors.UserAuthor, bool),
) *HardDeleteUserHandler {
	return &HardDeleteUserHandler{commandExecutor: commandExecutor, userId: userId, author: author}
}

func (h *HardDeleteUserHandler) Execute(c *gin.Context) (any, error) {
	userId, ok := h.userId(c.Request.Context())
	if !ok {
		return nil, app_errors.InternalError("MISSING_USER_ID")
	}
	author, ok := h.author(c.Request.Context())
	if !ok {
		return nil, app_errors.UnauthorizedError("MISSING_AUTHOR")
	}
	return nil, h.commandExecutor.SingleExecution(c.Request.Context(), userHardDelete.NewCommand(userId.String(), author))
}

func (h *HardDeleteUserHandler) Success(c *gin.Context, _ any) error {
	c.Status(http.StatusNoContent)
	return nil
}

func (h *HardDeleteUserHandler) Default(c *gin.Context) {
	c.Status(http.StatusInternalServerError)
}
