package scopes

import (
	"context"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/command"
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/protected_resources"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/scopes"
	scopeCreate "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/scopes/application/create"
)

type CreateScopeHandler struct {
	commandBus  command.Bus
	prOwnership func(context.Context) (protected_resources.Ownership, bool)
	prId        func(context.Context) (ids2.ProtectedResourceId, bool)
	scopeId     func(context.Context) (ids2.ScopeId, bool)
}

func NewCreateScopeHandler(
	commandBus command.Bus,
	prOwnership func(context.Context) (protected_resources.Ownership, bool),
	prId func(context.Context) (ids2.ProtectedResourceId, bool),
	scopeId func(context.Context) (ids2.ScopeId, bool),
) *CreateScopeHandler {
	return &CreateScopeHandler{commandBus: commandBus, prOwnership: prOwnership, prId: prId, scopeId: scopeId}
}

type createScopeRequest struct {
	Name string `json:"name" binding:"required"`
}

func (h *CreateScopeHandler) Execute(c *gin.Context) (any, error) {
	prOwnership, ok := h.prOwnership(c.Request.Context())
	if !ok {
		return nil, app_errors.InternalError("MISSING_PR_OWNERSHIP")
	}

	prId, ok := h.prId(c.Request.Context())
	if !ok {
		return nil, app_errors.InternalError("MISSING_PROTECTED_RESOURCE_ID")
	}

	scopeId, ok := h.scopeId(c.Request.Context())
	if !ok {
		return nil, app_errors.InternalError("MISSING_SCOPE_ID")
	}

	var req createScopeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, err
	}

	ownership := scopes.NewOwnershipFromVos(prId, prOwnership.UserId, prOwnership.AdminId)
	return nil, h.commandBus.Dispatch(c.Request.Context(), scopeCreate.NewCommand(scopeId, req.Name, ownership))
}

func (h *CreateScopeHandler) Success(c *gin.Context, _ any) error {
	c.Status(http.StatusCreated)
	return nil
}

func (h *CreateScopeHandler) Default(c *gin.Context) {
	c.Status(http.StatusInternalServerError)
}
