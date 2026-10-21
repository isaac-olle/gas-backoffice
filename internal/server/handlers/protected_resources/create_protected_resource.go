package protected_resources

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/protected_resources"
	prCreate "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/protected_resources/application/create"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/command"
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

type CreateProtectedResourceHandler struct {
	commandExecutor command.CommandExecutor
	prId            func(context.Context) (ids2.ProtectedResourceId, bool)
	subscriptionId  func(context.Context) (ids2.SubscriptionId, bool)
}

func NewCreateProtectedResourceHandler(
	commandExecutor command.CommandExecutor,
	prId func(context.Context) (ids2.ProtectedResourceId, bool),
	subscriptionId func(context.Context) (ids2.SubscriptionId, bool),
) *CreateProtectedResourceHandler {
	return &CreateProtectedResourceHandler{
		commandExecutor: commandExecutor,
		prId:            prId,
		subscriptionId:  subscriptionId,
	}
}

type createProtectedResourceRequest struct {
	Name string `json:"name" binding:"required"`
}

func (h *CreateProtectedResourceHandler) Execute(c *gin.Context) (any, error) {
	prId, ok := h.prId(c.Request.Context())
	if !ok {
		return nil, app_errors.InternalError("MISSING_PROTECTED_RESOURCE_ID")
	}

	subscriptionId, ok := h.subscriptionId(c.Request.Context())
	if !ok {
		return nil, app_errors.InternalError("MISSING_SUBSCRIPTION_ID")
	}

	var req createProtectedResourceRequest
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		return nil, app_errors.InvalidRequestBodyError(err)
	}

	prOwnership := protected_resources.NewOwnership(subscriptionId)
	return nil, h.commandExecutor.SingleExecution(c.Request.Context(), prCreate.NewCommand(prId, req.Name, prOwnership, false))
}

func (h *CreateProtectedResourceHandler) Success(c *gin.Context, _ any) error {
	c.Status(http.StatusCreated)
	return nil
}

func (h *CreateProtectedResourceHandler) Default(c *gin.Context) {
	c.Status(http.StatusInternalServerError)
}
