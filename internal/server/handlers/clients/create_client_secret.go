package clients

import (
	"context"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/command"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	clientSecretCreate "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/client_secrets/application/create"
)

type CreateClientSecretHandler struct {
	commandBus command.Bus
	clientId   func(context.Context) (ids.ClientId, bool)
}

func NewCreateClientSecretHandler(
	commandBus command.Bus,
	clientId func(context.Context) (ids.ClientId, bool),
) *CreateClientSecretHandler {
	return &CreateClientSecretHandler{commandBus: commandBus, clientId: clientId}
}

type createClientSecretRequest struct {
	SecretID        string `json:"secret_id" binding:"required"`
	PlaintextSecret string `json:"plaintext_secret" binding:"required"`
}

func (h *CreateClientSecretHandler) Execute(c *gin.Context) (any, error) {
	var req createClientSecretRequest
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		return nil, err
	}

	clientId, ok := h.clientId(c.Request.Context())
	if !ok {
		return nil, app_errors.InternalError("MISSING_CLIENT_ID")
	}

	return nil, h.commandBus.Dispatch(
		c.Request.Context(),
		clientSecretCreate.NewCommand(req.SecretID, clientId.String(), req.PlaintextSecret),
	)
}

func (h *CreateClientSecretHandler) Success(_ *gin.Context, _ any) error {
	return nil
}

func (h *CreateClientSecretHandler) Default(c *gin.Context) {
	c.Status(500)
}
