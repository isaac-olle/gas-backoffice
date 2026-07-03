package clients

import (
	"context"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/command"
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients"
	clientCreate "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients/application/create"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/tenants"
)

type CreateClientHandler struct {
	commandBus      command.Bus
	clientId        func(context.Context) (ids2.ClientId, bool)
	tenantId        func(context.Context) (ids2.TenantId, bool)
	tenantOwnership func(context.Context) (tenants.Ownership, bool)
}

func NewCreateClientHandler(
	commandBus command.Bus,
	clientId func(context.Context) (ids2.ClientId, bool),
	tenantId func(context.Context) (ids2.TenantId, bool),
	tenantOwnership func(context.Context) (tenants.Ownership, bool),
) *CreateClientHandler {
	return &CreateClientHandler{
		commandBus:      commandBus,
		clientId:        clientId,
		tenantId:        tenantId,
		tenantOwnership: tenantOwnership,
	}
}

type createClientRequest struct {
	Name string `json:"name" binding:"required"`
}

func (h *CreateClientHandler) Execute(c *gin.Context) (any, error) {
	var req createClientRequest
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		return nil, err
	}

	clientId, ok := h.clientId(c.Request.Context())
	if !ok {
		return nil, app_errors.InternalError("MISSING_CLIENT_ID")
	}

	tenantId, ok := h.tenantId(c.Request.Context())
	if !ok {
		return nil, app_errors.InternalError("MISSING_TENANT_ID")
	}

	ownership, ok := h.tenantOwnership(c.Request.Context())
	if !ok {
		return nil, app_errors.InternalError("MISSING_TENANT_OWNERSHIP")
	}

	return nil, h.commandBus.Dispatch(c.Request.Context(), clientCreate.NewCommand(
		clientId, req.Name, clients.NewOwnership(ownership.AdminId, tenantId), nil, nil,
	))
}

func (h *CreateClientHandler) Success(c *gin.Context, _ any) error {
	c.JSON(http.StatusCreated, gin.H{})
	return nil
}

func (h *CreateClientHandler) Default(c *gin.Context) {
	c.Status(http.StatusInternalServerError)
}
