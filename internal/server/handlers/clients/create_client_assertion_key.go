package clients

import (
	"context"
	"github.com/isaac-olle/gas/internal/shared/app_errors"
	"github.com/isaac-olle/gas/internal/shared/buses/command"
	ids2 "github.com/isaac-olle/gas/internal/shared/domain/vo/ids"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	clientAssertionKeyCreate "github.com/isaac-olle/gas/internal/bounded_contexts/security/client_assertion_keys/application/create"
)

type CreateClientAssertionKeyHandler struct {
	commandBus command.Bus
	clientId   func(context.Context) (ids2.ClientId, bool)
}

func NewCreateClientAssertionKeyHandler(
	commandBus command.Bus,
	clientId func(context.Context) (ids2.ClientId, bool),
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
		return nil, err
	}

	clientId, ok := h.clientId(c.Request.Context())
	if !ok {
		return nil, app_errors.InternalError("MISSING_CLIENT_ID")
	}

	keyId, err := ids2.NewClientAssertionKeyIdFromString(req.AssertionKeyID)
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
