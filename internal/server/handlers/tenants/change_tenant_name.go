package tenants

import (
	"context"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/command"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"net/http"

	"github.com/gin-gonic/gin"
	tenantChangeName "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/tenants/application/change_name"
)

type ChangeTenantNameHandler struct {
	commandBus command.Bus
	tenantId   func(context.Context) (ids.TenantId, bool)
}

func NewChangeTenantNameHandler(
	commandBus command.Bus,
	tenantId func(context.Context) (ids.TenantId, bool),
) *ChangeTenantNameHandler {
	return &ChangeTenantNameHandler{commandBus: commandBus, tenantId: tenantId}
}

type changeTenantNameRequest struct {
	Name string `json:"name" binding:"required"`
}

func (h *ChangeTenantNameHandler) Execute(c *gin.Context) (any, error) {
	tenantId, ok := h.tenantId(c.Request.Context())
	if !ok {
		return nil, app_errors.InternalError("MISSING_TENANT_ID")
	}

	var req changeTenantNameRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, err
	}

	return nil, h.commandBus.Dispatch(c.Request.Context(), tenantChangeName.NewCommand(tenantId, req.Name))
}

func (h *ChangeTenantNameHandler) Success(c *gin.Context, _ any) error {
	c.Status(http.StatusNoContent)
	return nil
}

func (h *ChangeTenantNameHandler) Default(c *gin.Context) {
	c.Status(http.StatusInternalServerError)
}
