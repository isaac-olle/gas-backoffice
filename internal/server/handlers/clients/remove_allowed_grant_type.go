package clients

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients/application/remove_allowed_grant_type"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/command"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

type RemoveAllowedGrantTypeHandler struct {
	commandExecutor command.CommandExecutor
	clientId        func(context.Context) (ids.ClientId, bool)
}

func NewRemoveAllowedGrantTypeHandler(
	commandExecutor command.CommandExecutor,
	clientId func(context.Context) (ids.ClientId, bool),
) *RemoveAllowedGrantTypeHandler {
	return &RemoveAllowedGrantTypeHandler{commandExecutor: commandExecutor, clientId: clientId}
}

type removeAllowedGrantTypeRequest struct {
	GrantType string `json:"grant_type" binding:"required"`
}

func (h *RemoveAllowedGrantTypeHandler) Execute(c *gin.Context) (any, error) {
	var req removeAllowedGrantTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, app_errors.InvalidRequestBodyError(err)
	}

	clientId, ok := h.clientId(c.Request.Context())
	if !ok {
		return nil, app_errors.InternalError("MISSING_CLIENT_ID")
	}

	return nil, h.commandExecutor.SingleExecution(c.Request.Context(), remove_allowed_grant_type.NewCommand(clientId, req.GrantType))
}

func (h *RemoveAllowedGrantTypeHandler) Success(c *gin.Context, _ any) error {
	c.Status(http.StatusNoContent)
	return nil
}

func (h *RemoveAllowedGrantTypeHandler) Default(c *gin.Context) {
	c.Status(http.StatusInternalServerError)
}
