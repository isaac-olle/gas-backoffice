package clients

import (
	"context"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/command"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"net/http"

	"github.com/gin-gonic/gin"
	clientSoftDelete "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients/application/soft_delete"
)

type SoftDeleteClientHandler struct {
	commandBus command.Bus
	clientId   func(context.Context) (ids.ClientId, bool)
}

func NewSoftDeleteClientHandler(
	commandBus command.Bus,
	clientId func(context.Context) (ids.ClientId, bool),
) *SoftDeleteClientHandler {
	return &SoftDeleteClientHandler{commandBus: commandBus, clientId: clientId}
}

func (h *SoftDeleteClientHandler) Execute(c *gin.Context) (any, error) {
	clientId, ok := h.clientId(c.Request.Context())
	if !ok {
		return nil, app_errors.InternalError("MISSING_CLIENT_ID")
	}
	return nil, h.commandBus.Dispatch(c.Request.Context(), clientSoftDelete.NewCommand(clientId.String()))
}

func (h *SoftDeleteClientHandler) Success(c *gin.Context, _ any) error {
	c.Status(http.StatusNoContent)
	return nil
}

func (h *SoftDeleteClientHandler) Default(c *gin.Context) {
	c.Status(http.StatusInternalServerError)
}
