package clients

import (
	"context"
	"github.com/isaac-olle/gas/internal/shared/app_errors"
	"github.com/isaac-olle/gas/internal/shared/buses/command"
	"github.com/isaac-olle/gas/internal/shared/domain/vo/ids"
	"net/http"

	"github.com/gin-gonic/gin"
	clientAddRedirectUri "github.com/isaac-olle/gas/internal/bounded_contexts/backoffice/clients/application/add_redirect_uri"
)

type AddRedirectUriHandler struct {
	commandBus command.Bus
	clientId   func(context.Context) (ids.ClientId, bool)
}

func NewAddRedirectUriHandler(
	commandBus command.Bus,
	clientId func(context.Context) (ids.ClientId, bool),
) *AddRedirectUriHandler {
	return &AddRedirectUriHandler{commandBus: commandBus, clientId: clientId}
}

type addRedirectUriRequest struct {
	RedirectUri string `json:"redirect_uri" binding:"required"`
}

func (h *AddRedirectUriHandler) Execute(c *gin.Context) (any, error) {
	var req addRedirectUriRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, err
	}

	clientId, ok := h.clientId(c.Request.Context())
	if !ok {
		return nil, app_errors.InternalError("MISSING_CLIENT_ID")
	}

	return nil, h.commandBus.Dispatch(c.Request.Context(), clientAddRedirectUri.NewCommand(clientId, req.RedirectUri))
}

func (h *AddRedirectUriHandler) Success(c *gin.Context, _ any) error {
	c.Status(http.StatusNoContent)
	return nil
}

func (h *AddRedirectUriHandler) Default(c *gin.Context) {
	c.Status(http.StatusInternalServerError)
}
