package create

import (
	"context"
	"time"

	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/command"
)

const CommandType command.Type = "command.backoffice.error_log.create"

type Command struct {
	errorId     string
	errorType   string
	code        string
	description string
	details     map[string]any
	cause       []string
	stackTrace  any
	requestId   string
	occurredAt  *time.Time
}

func NewCommand(
	errorId, errorType, code, description string,
	details map[string]any,
	cause []string,
	stackTrace any,
	requestId string,
	occurredAt *time.Time,
) *Command {
	return &Command{
		errorId:     errorId,
		errorType:   errorType,
		code:        code,
		description: description,
		details:     details,
		cause:       cause,
		stackTrace:  stackTrace,
		requestId:   requestId,
		occurredAt:  occurredAt,
	}
}

func (c *Command) Type() command.Type { return CommandType }

type CommandHandler struct{ service *Service }

func NewCommandHandler(service *Service) *CommandHandler { return &CommandHandler{service: service} }

func (h *CommandHandler) Handle(ctx context.Context, cmd command.Command) error {
	return command.Process[*Command, context.Context](ctx, cmd, h.service.Create)
}
