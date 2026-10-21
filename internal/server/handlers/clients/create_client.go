package clients

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients"
	clientCreate "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients/application/create"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/command"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

type CreateClientHandler struct {
	commandExecutor command.CommandExecutor
	clientId        func(context.Context) (ids.ClientId, bool)
	subscriptionId  func(context.Context) (ids.SubscriptionId, bool)
}

func NewCreateClientHandler(
	commandExecutor command.CommandExecutor,
	clientId func(context.Context) (ids.ClientId, bool),
	subscriptionId func(context.Context) (ids.SubscriptionId, bool),
) *CreateClientHandler {
	return &CreateClientHandler{
		commandExecutor: commandExecutor,
		clientId:        clientId,
		subscriptionId:  subscriptionId,
	}
}

type createClientRequest struct {
	Name              string   `json:"name" binding:"required"`
	AllowedGrantTypes []string `json:"allowed_grant_types" binding:"required,min=1"`
	RedirectUris      []string `json:"redirect_uris" binding:"required,min=1"`
	Public            bool     `json:"public"`
}

func (h *CreateClientHandler) Execute(c *gin.Context) (any, error) {
	var req createClientRequest
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		return nil, app_errors.InvalidRequestBodyError(err)
	}

	clientId, ok := h.clientId(c.Request.Context())
	if !ok {
		return nil, app_errors.InternalError("MISSING_CLIENT_ID")
	}

	subscriptionId, ok := h.subscriptionId(c.Request.Context())
	if !ok {
		return nil, app_errors.InternalError("MISSING_SUBSCRIPTION_ID")
	}

	return nil, h.commandExecutor.SingleExecution(c.Request.Context(), clientCreate.NewCommand(
		clientId, req.Name, clients.NewOwnership(subscriptionId), req.AllowedGrantTypes, req.RedirectUris, req.Public, false,
	))
}

func (h *CreateClientHandler) Success(c *gin.Context, _ any) error {
	c.JSON(http.StatusCreated, gin.H{})
	return nil
}

func (h *CreateClientHandler) Default(c *gin.Context) {
	c.Status(http.StatusInternalServerError)
}
