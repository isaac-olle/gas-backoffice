package clients

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients/application/remove_allowed_flow"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/command"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

type RemoveAllowedFlowHandler struct {
	commandBus command.Bus
	clientId   func(context.Context) (ids.ClientId, bool)
}

func NewRemoveAllowedFlowHandler(
	commandBus command.Bus,
	clientId func(context.Context) (ids.ClientId, bool),
) *RemoveAllowedFlowHandler {
	return &RemoveAllowedFlowHandler{commandBus: commandBus, clientId: clientId}
}

type removeAllowedFlowRequest struct {
	Flow string `json:"flow" binding:"required"`
}

func (h *RemoveAllowedFlowHandler) Execute(c *gin.Context) (any, error) {
	var req removeAllowedFlowRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, err
	}

	clientId, ok := h.clientId(c.Request.Context())
	if !ok {
		return nil, app_errors.InternalError("MISSING_CLIENT_ID")
	}

	return nil, h.commandBus.Dispatch(c.Request.Context(), remove_allowed_flow.NewCommand(clientId, req.Flow))
}

func (h *RemoveAllowedFlowHandler) Success(c *gin.Context, _ any) error {
	c.Status(http.StatusNoContent)
	return nil
}

func (h *RemoveAllowedFlowHandler) Default(c *gin.Context) {
	c.Status(http.StatusInternalServerError)
}
