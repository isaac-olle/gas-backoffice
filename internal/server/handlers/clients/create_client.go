package clients

import (
	"context"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/command"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients"
	clientCreate "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/clients/application/create"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/users"
)

type CreateClientHandler struct {
	commandBus    command.Bus
	clientId      func(context.Context) (ids.ClientId, bool)
	userId        func(context.Context) (ids.UserId, bool)
	userOwnership func(context.Context) (users.Ownership, bool)
}

func NewCreateClientHandler(
	commandBus command.Bus,
	clientId func(context.Context) (ids.ClientId, bool),
	userId func(context.Context) (ids.UserId, bool),
	userOwnership func(context.Context) (users.Ownership, bool),
) *CreateClientHandler {
	return &CreateClientHandler{
		commandBus:    commandBus,
		clientId:      clientId,
		userId:        userId,
		userOwnership: userOwnership,
	}
}

type createClientRequest struct {
	Name         string   `json:"name" binding:"required"`
	AllowedFlows []string `json:"allowed_flows" binding:"required,min=1"`
	RedirectUris []string `json:"redirect_uris" binding:"required,min=1"`
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

	userId, ok := h.userId(c.Request.Context())
	if !ok {
		return nil, app_errors.InternalError("MISSING_USER_ID")
	}

	ownership, ok := h.userOwnership(c.Request.Context())
	if !ok {
		return nil, app_errors.InternalError("MISSING_USER_OWNERSHIP")
	}

	return nil, h.commandBus.Dispatch(c.Request.Context(), clientCreate.NewCommand(
		clientId, req.Name, clients.NewOwnership(ownership.AdminId, userId), req.AllowedFlows, req.RedirectUris,
	))
}

func (h *CreateClientHandler) Success(c *gin.Context, _ any) error {
	c.JSON(http.StatusCreated, gin.H{})
	return nil
}

func (h *CreateClientHandler) Default(c *gin.Context) {
	c.Status(http.StatusInternalServerError)
}
