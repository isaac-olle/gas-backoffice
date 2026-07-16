package users

import (
	"context"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/command"
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/users"
	userCreate "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/users/application/create"
)

type CreateUserHandler struct {
	commandBus command.Bus
	adminId    func(context.Context) (ids2.AdminId, bool)
	userId     func(context.Context) (ids2.UserId, bool)
}

func NewCreateUserHandler(
	commandBus command.Bus,
	adminId func(context.Context) (ids2.AdminId, bool),
	userId func(context.Context) (ids2.UserId, bool),
) *CreateUserHandler {
	return &CreateUserHandler{commandBus: commandBus, adminId: adminId, userId: userId}
}

type createUserRequest struct {
	Name string `json:"name" binding:"required"`
}

func (h *CreateUserHandler) Execute(c *gin.Context) (any, error) {
	var req createUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, app_errors.InvalidRequestBodyError(err)
	}

	adminId, ok := h.adminId(c.Request.Context())
	if !ok {
		return nil, app_errors.InternalError("MISSING_ADMIN_ID")
	}

	userId, ok := h.userId(c.Request.Context())
	if !ok {
		return nil, app_errors.InternalError("MISSING_USER_ID")
	}

	return nil, h.commandBus.Dispatch(
		c.Request.Context(),
		userCreate.NewCommand(userId, req.Name, users.Ownership{AdminId: adminId}),
	)
}

func (h *CreateUserHandler) Success(c *gin.Context, _ any) error {
	c.Status(http.StatusCreated)
	return nil
}

func (h *CreateUserHandler) Default(c *gin.Context) {
	c.Status(http.StatusInternalServerError)
}
