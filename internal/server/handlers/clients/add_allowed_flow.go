package clients

import (
	"context"
	"github.com/isaac-olle/gas/internal/shared/app_errors"
	"github.com/isaac-olle/gas/internal/shared/buses/command"
	"github.com/isaac-olle/gas/internal/shared/domain/vo/ids"
	"net/http"

	"github.com/gin-gonic/gin"
	clientAddFlow "github.com/isaac-olle/gas/internal/bounded_contexts/backoffice/clients/application/add_allowed_flow"
)

type AddAllowedFlowHandler struct {
	commandBus command.Bus
	clientId   func(context.Context) (ids.ClientId, bool)
}

func NewAddAllowedFlowHandler(
	commandBus command.Bus,
	clientId func(context.Context) (ids.ClientId, bool),
) *AddAllowedFlowHandler {
	return &AddAllowedFlowHandler{commandBus: commandBus, clientId: clientId}
}

type addAllowedFlowRequest struct {
	Flow string `json:"flow" binding:"required"`
}

func (h *AddAllowedFlowHandler) Execute(c *gin.Context) (any, error) {
	var req addAllowedFlowRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, err
	}

	clientId, ok := h.clientId(c.Request.Context())
	if !ok {
		return nil, app_errors.InternalError("MISSING_CLIENT_ID")
	}

	return nil, h.commandBus.Dispatch(c.Request.Context(), clientAddFlow.NewCommand(clientId, req.Flow))
}

func (h *AddAllowedFlowHandler) Success(c *gin.Context, _ any) error {
	c.Status(http.StatusNoContent)
	return nil
}

func (h *AddAllowedFlowHandler) Default(c *gin.Context) {
	c.Status(http.StatusInternalServerError)
}
