package security

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	clientAssertionKeyCreate "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/client_assertion_keys/application/create"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/command"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

type CreateClientAssertionKeyHandler struct {
	commandBus command.Bus
	clientId   func(context.Context) (ids.ClientId, bool)
}

func NewCreateClientAssertionKeyHandler(
	commandBus command.Bus,
	clientId func(context.Context) (ids.ClientId, bool),
) *CreateClientAssertionKeyHandler {
	return &CreateClientAssertionKeyHandler{commandBus: commandBus, clientId: clientId}
}

type createClientAssertionKeyRequest struct {
	AssertionKeyID string `json:"assertion_key_id" binding:"required"`
	PublicKey      string `json:"public_key" binding:"required"`
}

func (h *CreateClientAssertionKeyHandler) Execute(c *gin.Context) (any, error) {
	var req createClientAssertionKeyRequest
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		return nil, app_errors.BadRequestError("INVALID_ASSERTION_KEY_REQUEST").WithCause(err)
	}

	clientId, ok := h.clientId(c.Request.Context())
	if !ok {
		return nil, app_errors.InternalError("MISSING_CLIENT_ID")
	}

	keyId, err := ids.NewClientAssertionKeyIdFromString(req.AssertionKeyID)
	if err != nil {
		return nil, app_errors.AddTrace(err)
	}

	return nil, h.commandBus.Dispatch(
		c.Request.Context(),
		clientAssertionKeyCreate.NewCommand(keyId, clientId, "RS256", req.PublicKey),
	)
}

func (h *CreateClientAssertionKeyHandler) Success(_ *gin.Context, _ any) error {
	return nil
}

func (h *CreateClientAssertionKeyHandler) Default(c *gin.Context) {
	c.Status(500)
}
