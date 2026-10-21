package subscriptions

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	addMember "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/subscriptions/application/add_member"
	createSubscription "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/subscriptions/application/create"
	removeMember "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/subscriptions/application/remove_member"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/command"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/authors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

type Handler struct {
	commands command.CommandExecutor
	author   func(context.Context) (*authors.UserAuthor, bool)
}

func NewHandler(commands command.CommandExecutor, author func(context.Context) (*authors.UserAuthor, bool)) *Handler {
	return &Handler{commands: commands, author: author}
}

type createRequest struct {
	Name string `json:"name" binding:"required"`
}

func (h *Handler) Execute(c *gin.Context) (any, error) {
	if c.Param("user_id") == "" {
		var req createRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			return nil, app_errors.InvalidRequestBodyError(err)
		}
		author, ok := h.author(c.Request.Context())
		if !ok {
			return nil, app_errors.UnauthorizedError("MISSING_AUTHOR")
		}
		id := ids.NewSubscriptionId()
		if err := h.commands.SingleExecution(c.Request.Context(), createSubscription.NewCommand(id, req.Name, author.Id)); err != nil {
			return nil, err
		}
		return gin.H{"id": id.String()}, nil
	}
	return h.changeMember(c, c.Request.Method == http.MethodPost)
}
func (h *Handler) changeMember(c *gin.Context, add bool) (any, error) {
	subscriptionId, err := ids.NewSubscriptionIdFromString(c.Param("subscription_id"))
	if err != nil {
		return nil, app_errors.BadRequestError("INVALID_SUBSCRIPTION_ID").WithCause(err, false)
	}
	userId, err := ids.NewUserIdFromString(c.Param("user_id"))
	if err != nil {
		return nil, app_errors.BadRequestError("INVALID_USER_ID").WithCause(err, false)
	}
	var cmd command.Command
	if add {
		cmd = addMember.NewCommand(subscriptionId, userId)
	} else {
		cmd = removeMember.NewCommand(subscriptionId, userId)
	}
	if err = h.commands.SingleExecution(c.Request.Context(), cmd); err != nil {
		return nil, err
	}
	return gin.H{"subscription_id": subscriptionId.String(), "user_id": userId.String()}, nil
}
func (h *Handler) Success(c *gin.Context, result any) error {
	if c.Request.Method == http.MethodPost && c.Param("user_id") == "" {
		c.JSON(http.StatusCreated, result)
	} else {
		c.JSON(http.StatusOK, result)
	}
	return nil
}
func (h *Handler) Default(c *gin.Context) { c.Status(http.StatusInternalServerError) }
