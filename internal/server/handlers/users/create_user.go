package users

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	userCreate "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/users/application/create"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/command"
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

type CreateUserHandler struct {
	commandBus command.Bus
	tenantId   func(context.Context) (ids2.UserId, bool)
	adminId    func(context.Context) (ids2.UserId, bool)
}

func NewCreateUserHandler(
	commandBus command.Bus,
	tenantId func(context.Context) (ids2.UserId, bool),
	adminId func(context.Context) (ids2.UserId, bool),
) *CreateUserHandler {
	return &CreateUserHandler{commandBus: commandBus, tenantId: tenantId, adminId: adminId}
}

type createUserRequest struct {
	Name string `json:"name" binding:"required"`
	Role string `json:"role" binding:"required"`
}

func (h *CreateUserHandler) Execute(c *gin.Context) (any, error) {
	var req createUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, app_errors.InvalidRequestBodyError(err)
	}

	tenantId, ok := h.tenantId(c.Request.Context())
	if !ok {
		return nil, app_errors.InternalError("MISSING_TENANT_ID")
	}

	adminId, ok := h.adminId(c.Request.Context())
	if !ok {
		return nil, app_errors.InternalError("MISSING_ADMIN_ID")
	}

	return nil, h.commandBus.Dispatch(
		c.Request.Context(),
		userCreate.NewCommand(tenantId, adminId, req.Name, req.Role),
	)
}

func (h *CreateUserHandler) Success(c *gin.Context, _ any) error {
	c.Status(http.StatusCreated)
	return nil
}

func (h *CreateUserHandler) Default(c *gin.Context) {
	c.Status(http.StatusInternalServerError)
}
