package subscriptions

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	updateSubscriptionName "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/subscriptions/application/update_name"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/command"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

type UpdateNameHandler struct {
	commands command.CommandExecutor
	author   func(context.Context) (*vo.UserAuthor, bool)
}

func NewUpdateNameHandler(commands command.CommandExecutor, author func(context.Context) (*vo.UserAuthor, bool)) *UpdateNameHandler {
	return &UpdateNameHandler{commands: commands, author: author}
}

type updateNameRequest struct {
	Name string `json:"name" binding:"required"`
}

func (h *UpdateNameHandler) Execute(c *gin.Context) (any, error) {
	subscriptionId, err := ids.NewSubscriptionIdFromString(c.Param("subscription_id"))
	if err != nil {
		return nil, app_errors.BadRequestError("INVALID_SUBSCRIPTION_ID").WithCause(err, false)
	}
	var req updateNameRequest
	if err = c.ShouldBindJSON(&req); err != nil {
		return nil, app_errors.InvalidRequestBodyError(err)
	}
	author, ok := h.author(c.Request.Context())
	if !ok {
		return nil, app_errors.UnauthorizedError("MISSING_AUTHOR")
	}
	return nil, h.commands.SingleExecution(c.Request.Context(), updateSubscriptionName.NewCommand(subscriptionId, req.Name, author.Id))
}

func (h *UpdateNameHandler) Success(c *gin.Context, _ any) error {
	c.Status(http.StatusNoContent)
	return nil
}

func (h *UpdateNameHandler) Default(c *gin.Context) { c.Status(http.StatusInternalServerError) }
