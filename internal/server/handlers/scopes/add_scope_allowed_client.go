package scopes

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	clientCheckOwnership "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients/application/check_ownership"
	scopeAddClient "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/scopes/application/add_allowed_client"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/command"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/query"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

type AddScopeAllowedClientHandler struct {
	commandExecutor command.CommandExecutor
	queryBus        query.Bus
	scopeId         func(context.Context) (ids.ScopeId, bool)
	clientId        func(context.Context) (ids.ClientId, bool)
	prId            func(context.Context) (ids.ProtectedResourceId, bool)
}

func NewAddScopeAllowedClientHandler(
	commandExecutor command.CommandExecutor,
	queryBus query.Bus,
	scopeId func(context.Context) (ids.ScopeId, bool),
	clientId func(context.Context) (ids.ClientId, bool),
	prId func(context.Context) (ids.ProtectedResourceId, bool),
) *AddScopeAllowedClientHandler {
	return &AddScopeAllowedClientHandler{commandExecutor: commandExecutor, queryBus: queryBus, scopeId: scopeId, clientId: clientId, prId: prId}
}

func (h *AddScopeAllowedClientHandler) Execute(c *gin.Context) (any, error) {
	scopeId, ok := h.scopeId(c.Request.Context())
	if !ok {
		return nil, app_errors.InternalError("MISSING_SCOPE_ID")
	}

	clientId, ok := h.clientId(c.Request.Context())
	if !ok {
		return nil, app_errors.InternalError("MISSING_CLIENT_ID")
	}

	prId, ok := h.prId(c.Request.Context())
	if !ok {
		return nil, app_errors.InternalError("MISSING_PROTECTED_RESOURCE_ID")
	}

	// Check the client belongs to the same subscription before linking it to the
	// scope: the scope's own ownership check only covers the scope/PR side, not the client.
	if _, err := h.queryBus.Dispatch(c.Request.Context(), clientCheckOwnership.NewQuery(clientId.String())); err != nil {
		return nil, err
	}

	return nil, h.commandExecutor.SingleExecution(
		c.Request.Context(), scopeAddClient.NewCommand(scopeId.String(), clientId.String(), prId.String()),
	)
}

func (h *AddScopeAllowedClientHandler) Success(c *gin.Context, _ any) error {
	c.Status(http.StatusNoContent)
	return nil
}

func (h *AddScopeAllowedClientHandler) Default(c *gin.Context) {
	c.Status(http.StatusInternalServerError)
}
