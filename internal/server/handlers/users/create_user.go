package users

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	createUser "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/users/application/create_user"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/command"
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

type CreateUserHandler struct {
	commandExecutor command.CommandExecutor
	userId          func(context.Context) (ids2.UserId, bool)
	authorId        func(context.Context) (ids2.UserId, bool)
}

func NewCreateUserHandler(
	commandExecutor command.CommandExecutor,
	userId func(context.Context) (ids2.UserId, bool),
	authorId func(context.Context) (ids2.UserId, bool),
) *CreateUserHandler {
	return &CreateUserHandler{commandExecutor: commandExecutor, userId: userId, authorId: authorId}
}

type createUserRequest struct {
	Name string `json:"name" binding:"required"`
}

func (h *CreateUserHandler) Execute(c *gin.Context) (any, error) {
	var req createUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, app_errors.InvalidRequestBodyError(err)
	}

	userId, ok := h.userId(c.Request.Context())
	if !ok {
		return nil, app_errors.InternalError("MISSING_USER_ID")
	}

	authorId, ok := h.authorId(c.Request.Context())
	if !ok {
		return nil, app_errors.InternalError("MISSING_AUTHOR_ID")
	}

	return nil, h.commandExecutor.SingleExecution(
		c.Request.Context(),
		createUser.NewCommand(userId, authorId, req.Name, false),
	)
}

func (h *CreateUserHandler) Success(c *gin.Context, _ any) error {
	c.Status(http.StatusCreated)
	return nil
}

func (h *CreateUserHandler) Default(c *gin.Context) {
	c.Status(http.StatusInternalServerError)
}
