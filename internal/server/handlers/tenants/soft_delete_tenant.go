package tenants

import (
	"context"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/command"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"net/http"

	"github.com/gin-gonic/gin"
	tenantSoftDelete "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/tenants/application/soft_delete"
)

type SoftDeleteTenantHandler struct {
	commandBus command.Bus
	tenantId   func(context.Context) (ids.TenantId, bool)
}

func NewSoftDeleteTenantHandler(
	commandBus command.Bus,
	tenantId func(context.Context) (ids.TenantId, bool),
) *SoftDeleteTenantHandler {
	return &SoftDeleteTenantHandler{commandBus: commandBus, tenantId: tenantId}
}

func (h *SoftDeleteTenantHandler) Execute(c *gin.Context) (any, error) {
	tenantId, ok := h.tenantId(c.Request.Context())
	if !ok {
		return nil, app_errors.InternalError("MISSING_TENANT_ID")
	}
	return nil, h.commandBus.Dispatch(c.Request.Context(), tenantSoftDelete.NewCommand(tenantId))
}

func (h *SoftDeleteTenantHandler) Success(c *gin.Context, _ any) error {
	c.Status(http.StatusNoContent)
	return nil
}

func (h *SoftDeleteTenantHandler) Default(c *gin.Context) {
	c.Status(http.StatusInternalServerError)
}
