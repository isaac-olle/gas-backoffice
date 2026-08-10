package sessions

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	sessionCreate "github.com/isaac-olle/gas-backoffice/internal/bounded_contexts/backoffice/auth/application/create"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/command"
	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/vo/ids"
)

type CreateSessionHandler struct {
	commandExecutor command.CommandExecutor
	sessionId       func(context.Context) (ids.SessionId, bool)
	userId          func(context.Context) (ids.UserId, bool)
}

func NewCreateSessionHandler(
	commandExecutor command.CommandExecutor,
	sessionId func(context.Context) (ids.SessionId, bool),
	userId func(context.Context) (ids.UserId, bool),
) *CreateSessionHandler {
	return &CreateSessionHandler{commandExecutor: commandExecutor, sessionId: sessionId, userId: userId}
}

func (h *CreateSessionHandler) Execute(c *gin.Context) (any, error) {
	sessionId, ok := h.sessionId(c.Request.Context())
	if !ok {
		return nil, app_errors.InternalError("MISSING_SESSION_ID")
	}

	userId, ok := h.userId(c.Request.Context())
	if !ok {
		return nil, app_errors.InternalError("MISSING_USER_ID")
	}

	return nil, h.commandExecutor.SingleExecution(
		c.Request.Context(),
		sessionCreate.NewCommand(sessionId, userId),
	)
}

func (h *CreateSessionHandler) Success(c *gin.Context, _ any) error {
	c.Status(http.StatusCreated)
	return nil
}

func (h *CreateSessionHandler) Default(c *gin.Context) {
	c.Status(http.StatusInternalServerError)
}
