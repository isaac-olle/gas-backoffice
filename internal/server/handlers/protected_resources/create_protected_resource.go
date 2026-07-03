package protected_resources

import (
	"context"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/command"
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/protected_resources"
	prCreate "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/protected_resources/application/create"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/tenants"
)

type CreateProtectedResourceHandler struct {
	commandBus      command.Bus
	tenantOwnership func(context.Context) (tenants.Ownership, bool)
	prId            func(context.Context) (ids2.ProtectedResourceId, bool)
	tenantId        func(context.Context) (ids2.TenantId, bool)
}

func NewCreateProtectedResourceHandler(
	commandBus command.Bus,
	tenantOwnership func(context.Context) (tenants.Ownership, bool),
	prId func(context.Context) (ids2.ProtectedResourceId, bool),
	tenantId func(context.Context) (ids2.TenantId, bool),
) *CreateProtectedResourceHandler {
	return &CreateProtectedResourceHandler{
		commandBus:      commandBus,
		tenantOwnership: tenantOwnership,
		prId:            prId,
		tenantId:        tenantId,
	}
}

type createProtectedResourceRequest struct {
	Name string `json:"name" binding:"required"`
}

func (h *CreateProtectedResourceHandler) Execute(c *gin.Context) (any, error) {
	tenantOwnership, ok := h.tenantOwnership(c.Request.Context())
	if !ok {
		return nil, app_errors.InternalError("MISSING_TENANT_OWNERSHIP")
	}

	prId, ok := h.prId(c.Request.Context())
	if !ok {
		return nil, app_errors.InternalError("MISSING_PROTECTED_RESOURCE_ID")
	}

	tenantId, ok := h.tenantId(c.Request.Context())
	if !ok {
		return nil, app_errors.InternalError("MISSING_TENANT_ID")
	}

	var req createProtectedResourceRequest
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		return nil, err
	}

	prOwnership := protected_resources.NewOwnership(tenantId, tenantOwnership.AdminId)
	return nil, h.commandBus.Dispatch(c.Request.Context(), prCreate.NewCommand(prId, req.Name, prOwnership))
}

func (h *CreateProtectedResourceHandler) Success(c *gin.Context, _ any) error {
	c.Status(http.StatusCreated)
	return nil
}

func (h *CreateProtectedResourceHandler) Default(c *gin.Context) {
	c.Status(http.StatusInternalServerError)
}
