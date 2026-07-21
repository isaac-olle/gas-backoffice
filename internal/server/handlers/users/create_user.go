package users

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/users"
	userCreate "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/users/application/create"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/command"
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

type CreateUserHandler struct {
	commandBus command.Bus
	userId     func(context.Context) (ids2.UserId, bool)
}

func NewCreateUserHandler(
	commandBus command.Bus,
	userId func(context.Context) (ids2.UserId, bool),
) *CreateUserHandler {
	return &CreateUserHandler{commandBus: commandBus, userId: userId}
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

	return nil, h.commandBus.Dispatch(
		c.Request.Context(),
		userCreate.NewCommand(userId, req.Name, users.NewOwnership(userId)),
	)
}

func (h *CreateUserHandler) Success(c *gin.Context, _ any) error {
	c.Status(http.StatusCreated)
	return nil
}

func (h *CreateUserHandler) Default(c *gin.Context) {
	c.Status(http.StatusInternalServerError)
}
