package users

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	createTenant "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/tenants/application/create_tenant"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/command"
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

type CreateUserHandler struct {
	commandExecutor command.CommandExecutor
	tenantId        func(context.Context) (ids2.UserId, bool)
	adminId         func(context.Context) (ids2.UserId, bool)
}

func NewCreateUserHandler(
	commandExecutor command.CommandExecutor,
	tenantId func(context.Context) (ids2.UserId, bool),
	adminId func(context.Context) (ids2.UserId, bool),
) *CreateUserHandler {
	return &CreateUserHandler{commandExecutor: commandExecutor, tenantId: tenantId, adminId: adminId}
}

type createUserRequest struct {
	Name string `json:"name" binding:"required"`
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

	return nil, h.commandExecutor.SingleExecution(
		c.Request.Context(),
		createTenant.NewCommand(tenantId, adminId, req.Name, false),
	)
}

func (h *CreateUserHandler) Success(c *gin.Context, _ any) error {
	c.Status(http.StatusCreated)
	return nil
}

func (h *CreateUserHandler) Default(c *gin.Context) {
	c.Status(http.StatusInternalServerError)
}
