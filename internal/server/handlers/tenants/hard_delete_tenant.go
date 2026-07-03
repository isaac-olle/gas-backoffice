package tenants

import (
	"context"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/command"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"net/http"

	"github.com/gin-gonic/gin"
	tenantHardDelete "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/tenants/application/hard_delete"
)

type HardDeleteTenantHandler struct {
	commandBus command.Bus
	tenantId   func(context.Context) (ids.TenantId, bool)
}

func NewHardDeleteTenantHandler(
	commandBus command.Bus,
	tenantId func(context.Context) (ids.TenantId, bool),
) *HardDeleteTenantHandler {
	return &HardDeleteTenantHandler{commandBus: commandBus, tenantId: tenantId}
}

func (h *HardDeleteTenantHandler) Execute(c *gin.Context) (any, error) {
	tenantId, ok := h.tenantId(c.Request.Context())
	if !ok {
		return nil, app_errors.InternalError("MISSING_TENANT_ID")
	}
	return nil, h.commandBus.Dispatch(c.Request.Context(), tenantHardDelete.NewCommand(tenantId.String()))
}

func (h *HardDeleteTenantHandler) Success(c *gin.Context, _ any) error {
	c.Status(http.StatusNoContent)
	return nil
}

func (h *HardDeleteTenantHandler) Default(c *gin.Context) {
	c.Status(http.StatusInternalServerError)
}
