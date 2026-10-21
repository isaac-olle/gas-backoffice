package clients

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	clientAddGrantType "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients/application/add_allowed_grant_type"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/command"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

type AddAllowedGrantTypeHandler struct {
	commandExecutor command.CommandExecutor
	clientId        func(context.Context) (ids.ClientId, bool)
}

func NewAddAllowedGrantTypeHandler(
	commandExecutor command.CommandExecutor,
	clientId func(context.Context) (ids.ClientId, bool),
) *AddAllowedGrantTypeHandler {
	return &AddAllowedGrantTypeHandler{commandExecutor: commandExecutor, clientId: clientId}
}

type addAllowedGrantTypeRequest struct {
	GrantType string `json:"grant_type" binding:"required"`
}

func (h *AddAllowedGrantTypeHandler) Execute(c *gin.Context) (any, error) {
	var req addAllowedGrantTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, app_errors.InvalidRequestBodyError(err)
	}

	clientId, ok := h.clientId(c.Request.Context())
	if !ok {
		return nil, app_errors.InternalError("MISSING_CLIENT_ID")
	}

	return nil, h.commandExecutor.SingleExecution(c.Request.Context(), clientAddGrantType.NewCommand(clientId, req.GrantType))
}

func (h *AddAllowedGrantTypeHandler) Success(c *gin.Context, _ any) error {
	c.Status(http.StatusNoContent)
	return nil
}

func (h *AddAllowedGrantTypeHandler) Default(c *gin.Context) {
	c.Status(http.StatusInternalServerError)
}
