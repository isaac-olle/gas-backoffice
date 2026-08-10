package security

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	clientAssertionKeyCreate "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/security/client_assertion_keys/application/create"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/command"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

type CreateClientAssertionKeyHandler struct {
	commandExecutor command.CommandExecutor
	clientId        func(context.Context) (ids.ClientId, bool)
	assertionKeyId  func(context.Context) (ids.ClientAssertionKeyId, bool)
}

func NewCreateClientAssertionKeyHandler(
	commandExecutor command.CommandExecutor,
	clientId func(context.Context) (ids.ClientId, bool),
	assertionKeyId func(context.Context) (ids.ClientAssertionKeyId, bool),
) *CreateClientAssertionKeyHandler {
	return &CreateClientAssertionKeyHandler{
		commandExecutor: commandExecutor,
		clientId:        clientId,
		assertionKeyId:  assertionKeyId,
	}
}

type createClientAssertionKeyRequest struct {
	PublicKey string `json:"public_key" binding:"required"`
}

func (h *CreateClientAssertionKeyHandler) Execute(c *gin.Context) (any, error) {
	var req createClientAssertionKeyRequest
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		return nil, app_errors.BadRequestError("INVALID_ASSERTION_KEY_REQUEST").WithDescription(err.Error())
	}

	clientId, ok := h.clientId(c.Request.Context())
	if !ok {
		return nil, app_errors.InternalError("MISSING_CLIENT_ID")
	}

	keyId, ok := h.assertionKeyId(c.Request.Context())
	if !ok {
		return nil, app_errors.InternalError("MISSING_ASSERTION_KEY_ID")
	}

	return nil, h.commandExecutor.SingleExecution(
		c.Request.Context(),
		clientAssertionKeyCreate.NewCommand(keyId, clientId, "RS256", req.PublicKey),
	)
}

func (h *CreateClientAssertionKeyHandler) Success(c *gin.Context, _ any) error {
	c.Status(http.StatusCreated)
	return nil
}

func (h *CreateClientAssertionKeyHandler) Default(c *gin.Context) {
	c.Status(http.StatusInternalServerError)
}
