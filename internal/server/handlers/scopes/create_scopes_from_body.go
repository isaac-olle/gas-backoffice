package scopes

import (
	"context"
	"github.com/isaac-olle/gas/internal/shared/app_errors"
	"github.com/isaac-olle/gas/internal/shared/buses/command"
	ids2 "github.com/isaac-olle/gas/internal/shared/domain/vo/ids"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/isaac-olle/gas/internal/bounded_contexts/backoffice/scopes"
	scopeCreate "github.com/isaac-olle/gas/internal/bounded_contexts/backoffice/scopes/application/create"
	"github.com/isaac-olle/gas/internal/bounded_contexts/backoffice/tenants"
)

type CreateScopesFromBodyHandler struct {
	commandBus      command.Bus
	tenantOwnership func(context.Context) (tenants.Ownership, bool)
	prId            func(context.Context) (ids2.ProtectedResourceId, bool)
	tenantId        func(context.Context) (ids2.TenantId, bool)
	scopeId         func(context.Context) (ids2.ScopeId, bool)
}

func NewCreateScopesFromBodyHandler(
	commandBus command.Bus,
	tenantOwnership func(context.Context) (tenants.Ownership, bool),
	prId func(context.Context) (ids2.ProtectedResourceId, bool),
	tenantId func(context.Context) (ids2.TenantId, bool),
	scopeId func(context.Context) (ids2.ScopeId, bool),
) *CreateScopesFromBodyHandler {
	return &CreateScopesFromBodyHandler{
		commandBus:      commandBus,
		tenantOwnership: tenantOwnership,
		prId:            prId,
		tenantId:        tenantId,
		scopeId:         scopeId,
	}
}

type createScopesFromBodyRequest struct {
	Scopes []string `json:"scopes"`
}

func (h *CreateScopesFromBodyHandler) Execute(c *gin.Context) (any, error) {
	var req createScopesFromBodyRequest
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		return nil, err
	}

	if len(req.Scopes) == 0 {
		return nil, nil
	}

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

	for _, scopeName := range req.Scopes {
		scopeId, _ := h.scopeId(c.Request.Context())
		ownership := scopes.NewOwnershipFromVos(prId, tenantId, tenantOwnership.AdminId)
		if err := h.commandBus.Dispatch(
			c.Request.Context(), scopeCreate.NewCommand(scopeId, scopeName, ownership),
		); err != nil {
			return nil, err
		}
	}

	return nil, nil
}

func (h *CreateScopesFromBodyHandler) Success(_ *gin.Context, _ any) error {
	return nil
}

func (h *CreateScopesFromBodyHandler) Default(c *gin.Context) {
	c.Status(500)
}
