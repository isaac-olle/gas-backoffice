package scopes

import (
	"context"
	prCheckOwnership "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/protected_resources/application/check_ownership"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/scopes"
	scopeCreate "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/scopes/application/create"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/command"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/query"
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type CreateScopesFromBodyHandler struct {
	commandExecutor command.CommandExecutor
	queryBus        query.Bus
	prId            func(context.Context) (ids2.ProtectedResourceId, bool)
	scopeId         func(context.Context) (ids2.ScopeId, bool)
}

func NewCreateScopesFromBodyHandler(
	commandExecutor command.CommandExecutor,
	queryBus query.Bus,
	prId func(context.Context) (ids2.ProtectedResourceId, bool),
	scopeId func(context.Context) (ids2.ScopeId, bool),
) *CreateScopesFromBodyHandler {
	return &CreateScopesFromBodyHandler{
		commandExecutor: commandExecutor,
		queryBus:        queryBus,
		prId:            prId,
		scopeId:         scopeId,
	}
}

type createScopesFromBodyRequest struct {
	Scopes []string `json:"scopes"`
}

func (h *CreateScopesFromBodyHandler) Execute(c *gin.Context) (any, error) {
	var req createScopesFromBodyRequest
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		return nil, app_errors.InvalidRequestBodyError(err)
	}

	if len(req.Scopes) == 0 {
		return nil, nil
	}

	prId, ok := h.prId(c.Request.Context())
	if !ok {
		return nil, app_errors.InternalError("MISSING_PROTECTED_RESOURCE_ID")
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

	for _, scopeName := range req.Scopes {
		scopeId, _ := h.scopeId(c.Request.Context())
		if err := h.commandExecutor.SingleExecution(
			c.Request.Context(), scopeCreate.NewCommand(scopeId, scopeName, ownership, false),
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
