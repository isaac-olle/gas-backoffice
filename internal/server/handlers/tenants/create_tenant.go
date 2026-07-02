package tenants

import (
	"context"
	"github.com/isaac-olle/gas/internal/shared/app_errors"
	"github.com/isaac-olle/gas/internal/shared/buses/command"
	ids2 "github.com/isaac-olle/gas/internal/shared/domain/vo/ids"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/isaac-olle/gas/internal/bounded_contexts/backoffice/tenants"
	tenantCreate "github.com/isaac-olle/gas/internal/bounded_contexts/backoffice/tenants/application/create"
)

type CreateTenantHandler struct {
	commandBus command.Bus
	adminId    func(context.Context) (ids2.AdminId, bool)
	tenantId   func(context.Context) (ids2.TenantId, bool)
}

func NewCreateTenantHandler(
	commandBus command.Bus,
	adminId func(context.Context) (ids2.AdminId, bool),
	tenantId func(context.Context) (ids2.TenantId, bool),
) *CreateTenantHandler {
	return &CreateTenantHandler{commandBus: commandBus, adminId: adminId, tenantId: tenantId}
}

type createTenantRequest struct {
	Name string `json:"name" binding:"required"`
}

func (h *CreateTenantHandler) Execute(c *gin.Context) (any, error) {
	var req createTenantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, err
	}

	adminId, ok := h.adminId(c.Request.Context())
	if !ok {
		return nil, app_errors.InternalError("MISSING_ADMIN_ID")
	}

	tenantId, ok := h.tenantId(c.Request.Context())
	if !ok {
		return nil, app_errors.InternalError("MISSING_TENANT_ID")
	}

	return nil, h.commandBus.Dispatch(
		c.Request.Context(),
		tenantCreate.NewCommand(tenantId, req.Name, tenants.Ownership{AdminId: adminId}),
	)
}

func (h *CreateTenantHandler) Success(c *gin.Context, _ any) error {
	c.Status(http.StatusCreated)
	return nil
}

func (h *CreateTenantHandler) Default(c *gin.Context) {
	c.Status(http.StatusInternalServerError)
}
