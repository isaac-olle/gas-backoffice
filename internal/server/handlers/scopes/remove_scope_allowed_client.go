package scopes

import (
	"context"
	"github.com/isaac-olle/gas/internal/shared/app_errors"
	"github.com/isaac-olle/gas/internal/shared/buses/command"
	ids2 "github.com/isaac-olle/gas/internal/shared/domain/vo/ids"
	"net/http"

	"github.com/gin-gonic/gin"
	scopeRemoveClient "github.com/isaac-olle/gas/internal/bounded_contexts/backoffice/scopes/application/remove_allowed_client"
)

type RemoveScopeAllowedClientHandler struct {
	commandBus command.Bus
	scopeId    func(context.Context) (ids2.ScopeId, bool)
	clientId   func(context.Context) (ids2.ClientId, bool)
}

func NewRemoveScopeAllowedClientHandler(
	commandBus command.Bus,
	scopeId func(context.Context) (ids2.ScopeId, bool),
	clientId func(context.Context) (ids2.ClientId, bool),
) *RemoveScopeAllowedClientHandler {
	return &RemoveScopeAllowedClientHandler{commandBus: commandBus, scopeId: scopeId, clientId: clientId}
}

func (h *RemoveScopeAllowedClientHandler) Execute(c *gin.Context) (any, error) {
	scopeId, ok := h.scopeId(c.Request.Context())
	if !ok {
		return nil, app_errors.InternalError("MISSING_SCOPE_ID")
	}

	clientId, ok := h.clientId(c.Request.Context())
	if !ok {
		return nil, app_errors.InternalError("MISSING_CLIENT_ID")
	}

	return nil, h.commandBus.Dispatch(
		c.Request.Context(), scopeRemoveClient.NewCommand(scopeId.String(), clientId.String()),
	)
}

func (h *RemoveScopeAllowedClientHandler) Success(c *gin.Context, _ any) error {
	c.Status(http.StatusNoContent)
	return nil
}

func (h *RemoveScopeAllowedClientHandler) Default(c *gin.Context) {
	c.Status(http.StatusInternalServerError)
}
