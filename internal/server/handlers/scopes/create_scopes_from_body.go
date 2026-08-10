package scopes

import (
	"context"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/command"
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/scopes"
	scopeCreate "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/scopes/application/create"
)

type CreateScopesFromBodyHandler struct {
	commandExecutor command.CommandExecutor
	prId            func(context.Context) (ids2.ProtectedResourceId, bool)
	userId          func(context.Context) (ids2.UserId, bool)
	scopeId         func(context.Context) (ids2.ScopeId, bool)
}

func NewCreateScopesFromBodyHandler(
	commandExecutor command.CommandExecutor,
	prId func(context.Context) (ids2.ProtectedResourceId, bool),
	userId func(context.Context) (ids2.UserId, bool),
	scopeId func(context.Context) (ids2.ScopeId, bool),
) *CreateScopesFromBodyHandler {
	return &CreateScopesFromBodyHandler{
		commandExecutor: commandExecutor,
		prId:            prId,
		userId:          userId,
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

	userId, ok := h.userId(c.Request.Context())
	if !ok {
		return nil, app_errors.InternalError("MISSING_USER_ID")
	}

	for _, scopeName := range req.Scopes {
		scopeId, _ := h.scopeId(c.Request.Context())
		ownership := scopes.NewOwnershipFromVos(prId, userId)
		if err := h.commandExecutor.SingleExecution(
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
