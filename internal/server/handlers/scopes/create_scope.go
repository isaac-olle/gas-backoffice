package scopes

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	prCheckOwnership "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/protected_resources/application/check_ownership"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/scopes"
	scopeCreate "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/scopes/application/create"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/command"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/query"
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

type CreateScopeHandler struct {
	commandExecutor command.CommandExecutor
	queryBus        query.Bus
	prId            func(context.Context) (ids2.ProtectedResourceId, bool)
	scopeId         func(context.Context) (ids2.ScopeId, bool)
}

func NewCreateScopeHandler(
	commandExecutor command.CommandExecutor,
	queryBus query.Bus,
	prId func(context.Context) (ids2.ProtectedResourceId, bool),
	scopeId func(context.Context) (ids2.ScopeId, bool),
) *CreateScopeHandler {
	return &CreateScopeHandler{commandExecutor: commandExecutor, queryBus: queryBus, prId: prId, scopeId: scopeId}
}

type createScopeRequest struct {
	Name string `json:"name" binding:"required"`
}

func (h *CreateScopeHandler) Execute(c *gin.Context) (any, error) {
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
		return nil, app_errors.InvalidRequestBodyError(err)
	}

	result, err := h.queryBus.Dispatch(c.Request.Context(), prCheckOwnership.NewQuery(prId.String()))
	if err != nil {
		return nil, err
	}
	subscriptionIdString, ok := result.(string)
	if !ok {
		return nil, app_errors.InternalError("INVALID_PROTECTED_RESOURCE_OWNERSHIP")
	}
	subscriptionId, err := ids2.NewSubscriptionIdFromString(subscriptionIdString)
	if err != nil {
		return nil, app_errors.InternalError("INVALID_PROTECTED_RESOURCE_OWNERSHIP").WithCause(err, false)
	}
	ownership := scopes.NewOwnershipFromVos(prId, subscriptionId)
	return nil, h.commandExecutor.SingleExecution(c.Request.Context(), scopeCreate.NewCommand(scopeId, req.Name, ownership, false))
}

func (h *CreateScopeHandler) Success(c *gin.Context, _ any) error {
	c.Status(http.StatusCreated)
	return nil
}

func (h *CreateScopeHandler) Default(c *gin.Context) {
	c.Status(http.StatusInternalServerError)
}
