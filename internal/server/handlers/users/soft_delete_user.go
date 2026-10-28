package users

import (
	"context"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/command"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/authors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"net/http"

	"github.com/gin-gonic/gin"
	userSoftDelete "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/users/application/soft_delete"
)

type SoftDeleteUserHandler struct {
	commandExecutor command.CommandExecutor
	userId          func(context.Context) (ids.UserId, bool)
	author          func(context.Context) (*authors.UserAuthor, bool)
}

func NewSoftDeleteUserHandler(
	commandExecutor command.CommandExecutor,
	userId func(context.Context) (ids.UserId, bool),
	author func(context.Context) (*authors.UserAuthor, bool),
) *SoftDeleteUserHandler {
	return &SoftDeleteUserHandler{commandExecutor: commandExecutor, userId: userId, author: author}
}

func (h *SoftDeleteUserHandler) Execute(c *gin.Context) (any, error) {
	userId, ok := h.userId(c.Request.Context())
	if !ok {
		return nil, app_errors.InternalError("MISSING_USER_ID")
	}
	author, ok := h.author(c.Request.Context())
	if !ok {
		return nil, app_errors.UnauthorizedError("MISSING_AUTHOR")
	}
	return nil, h.commandExecutor.SingleExecution(c.Request.Context(), userSoftDelete.NewCommand(userId, author))
}

func (h *SoftDeleteUserHandler) Success(c *gin.Context, _ any) error {
	c.Status(http.StatusNoContent)
	return nil
}

func (h *SoftDeleteUserHandler) Default(c *gin.Context) {
	c.Status(http.StatusInternalServerError)
}
