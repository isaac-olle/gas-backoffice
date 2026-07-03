package tenants

import (
	"context"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/query"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/tenants"
	tenantGet "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/tenants/application/get"
)

type GetTenantHandler struct {
	queryBus query.Bus
	tenantId func(context.Context) (ids.TenantId, bool)
}

func NewGetTenantHandler(queryBus query.Bus, tenantId func(context.Context) (ids.TenantId, bool)) *GetTenantHandler {
	return &GetTenantHandler{queryBus: queryBus, tenantId: tenantId}
}

func (h *GetTenantHandler) Execute(c *gin.Context) (any, error) {
	tenantId, ok := h.tenantId(c.Request.Context())
	if !ok {
		return nil, app_errors.InternalError("MISSING_TENANT_ID")
	}
	tenant, err := query.DispatchTyped[*tenants.Tenant](c.Request.Context(), h.queryBus, tenantGet.NewQuery(tenantId))
	if err != nil {
		return nil, err
	}
	return tenant.ToPrimitives(), nil
}

func (h *GetTenantHandler) Success(c *gin.Context, result any) error {
	c.JSON(http.StatusOK, result)
	return nil
}

func (h *GetTenantHandler) Default(c *gin.Context) {
	c.Status(http.StatusInternalServerError)
}
