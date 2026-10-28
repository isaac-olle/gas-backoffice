package users

import (
	"context"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/command"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/authors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"net/http"

	"github.com/gin-gonic/gin"
	userChangeName "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/users/application/change_name"
)

type ChangeUserNameHandler struct {
	commandExecutor command.CommandExecutor
	userId          func(context.Context) (ids.UserId, bool)
	author          func(context.Context) (*authors.UserAuthor, bool)
}

func NewChangeUserNameHandler(
	commandExecutor command.CommandExecutor,
	userId func(context.Context) (ids.UserId, bool),
	author func(context.Context) (*authors.UserAuthor, bool),
) *ChangeUserNameHandler {
	return &ChangeUserNameHandler{commandExecutor: commandExecutor, userId: userId, author: author}
}

type changeUserNameRequest struct {
	Name string `json:"name" binding:"required"`
}

func (h *ChangeUserNameHandler) Execute(c *gin.Context) (any, error) {
	userId, ok := h.userId(c.Request.Context())
	if !ok {
		return nil, app_errors.InternalError("MISSING_USER_ID")
	}
	author, ok := h.author(c.Request.Context())
	if !ok {
		return nil, app_errors.UnauthorizedError("MISSING_AUTHOR")
	}

	var req changeUserNameRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, app_errors.InvalidRequestBodyError(err)
	}

	return nil, h.commandExecutor.SingleExecution(c.Request.Context(), userChangeName.NewCommand(userId, req.Name, author))
}

func (h *ChangeUserNameHandler) Success(c *gin.Context, _ any) error {
	c.Status(http.StatusNoContent)
	return nil
}

func (h *ChangeUserNameHandler) Default(c *gin.Context) {
	c.Status(http.StatusInternalServerError)
}
