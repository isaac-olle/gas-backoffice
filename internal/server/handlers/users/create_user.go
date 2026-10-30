package users

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	addMember "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/subscriptions/application/add_member"
	createUser "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/users/application/create_user"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/command"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo"
	ids2 "github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

type CreateUserHandler struct {
	commandExecutor command.CommandExecutor
	userId          func(context.Context) (ids2.UserId, bool)
	authorId        func(context.Context) (ids2.UserId, bool)
}

func NewCreateUserHandler(
	commandExecutor command.CommandExecutor,
	userId func(context.Context) (ids2.UserId, bool),
	authorId func(context.Context) (ids2.UserId, bool),
) *CreateUserHandler {
	return &CreateUserHandler{commandExecutor: commandExecutor, userId: userId, authorId: authorId}
}

type createUserRequest struct {
	Name            string   `json:"name" binding:"required"`
	Email           string   `json:"email" binding:"required,email"`
	Role            string   `json:"role" binding:"required"`
	IsSystem        bool     `json:"is_system"`
	SubscriptionIds []string `json:"subscription_ids"`
}

func (h *CreateUserHandler) Execute(c *gin.Context) (any, error) {
	var req createUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, app_errors.InvalidRequestBodyError(err)
	}
	role, err := vo.NewRole(req.Role)
	if err != nil || role != vo.SubscriptionUserRole {
		return nil, app_errors.BadRequestError("INVALID_USER_ROLE").
			WithDescription(fmt.Sprintf("unsupported user role %q", req.Role))
	}

	userId, ok := h.userId(c.Request.Context())
	if !ok {
		return nil, app_errors.InternalError("MISSING_USER_ID")
	}

	authorId, ok := h.authorId(c.Request.Context())
	if !ok {
		return nil, app_errors.InternalError("MISSING_AUTHOR_ID")
	}

	subscriptionIds := make([]ids2.SubscriptionId, 0, len(req.SubscriptionIds))
	for _, rawSubscriptionId := range req.SubscriptionIds {
		subscriptionId, err := ids2.NewSubscriptionIdFromString(rawSubscriptionId)
		if err != nil {
			return nil, app_errors.BadRequestError("INVALID_SUBSCRIPTION_ID").WithCause(err, false)
		}
		subscriptionIds = append(subscriptionIds, subscriptionId)
	}

	commands := []command.Command{
		createUser.NewCommand(userId, authorId, req.Name, role, req.IsSystem, req.Email),
	}
	for _, subscriptionId := range subscriptionIds {
		commands = append(commands, addMember.NewCommand(subscriptionId, userId, vo.UserSubscriptionRole))
	}

	return nil, h.commandExecutor.MultipleExecution(c.Request.Context(), commands...)
}

func (h *CreateUserHandler) Success(c *gin.Context, _ any) error {
	c.Status(http.StatusCreated)
	return nil
}

func (h *CreateUserHandler) Default(c *gin.Context) {
	c.Status(http.StatusInternalServerError)
}
